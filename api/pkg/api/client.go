package api

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/cenkalti/backoff/v4"
)

const (
	MethodPost = http.MethodPost
)

const (
	StatusOK = http.StatusOK
)

const (
	ContentTypeApplicationJSON = "application/json"
)

type Client struct {
	// HTTP client to be used to execute HTTP call.
	HTTPClient *http.Client

	// URL for this Client. This value can contain path
	// variable placeholders for substitution when sending the HTTP request
	url string

	// HTTP method
	method string

	// Name of the service to call. Together with serviceName will form the label in logger fields & error code prefixes.
	serviceName string

	// Name of the call. Together with serviceName will form the label in logger fields & error code prefixes.
	callName string

	// Content MIME type
	// Default: application/json
	contentType string

	// Default request header configuration
	header header

	timeoutAndRetryOption timeoutAndRetryOption
}

// timeoutAndRetryOption holds timeout & retry info for client. This is optional. If not provided, it will
// pick up the default config.
type timeoutAndRetryOption struct {
	// Max num of retries. Setting to <= 0 means no retry
	// Default: 0
	maxRetries uint64
	// Max execution wait time per try.
	// Default: 15s
	maxWaitPerTry time.Duration
	// Max execution wait time, regardless of retries.
	// Default: 15s
	maxWaitInclRetries time.Duration
	// Set false to exclude retry on timeout errors.
	// Good for non-idempotent resources (i.e. push notifications)
	// Default: false
	onTimeout bool
	// Retry on certain http status code
	// Default: empty
	onStatusCodes map[int]bool
}

// header describes the HTTP request headers
type header struct {
	// Other request header values
	// Note: Can be overridden on Payload
	// Default: nil
	values map[string]string
}

// NewClient method creates a new client
func NewClient(
	url,
	method,
	callName string,
	contentType string,
) (*Client, error) {
	c := &Client{
		HTTPClient: &http.Client{Timeout: 1 * time.Minute},
		timeoutAndRetryOption: timeoutAndRetryOption{
			maxRetries:         3,
			maxWaitPerTry:      1 * time.Minute,
			maxWaitInclRetries: 1 * time.Minute,
			onTimeout:          true,
			onStatusCodes:      make(map[int]bool),
		},
		contentType: contentType,
	}

	c.url = strings.TrimSpace(url)
	if c.url == "" {
		return nil, errors.New("url is missing")
	}
	c.method = strings.TrimSpace(method)
	if c.method == "" {
		return nil, errors.New("method is missing")
	}

	c.callName = strings.TrimSpace(callName)
	if c.callName == "" {
		return nil, errors.New("callName is missing")
	}

	return c, nil
}

// constructURL returns the full URL with query params and path variables substitution
func (c *Client) constructURL(p Payload) string {
	u := c.url
	// Replace path variables
	for k, v := range p.PathVars {
		u = strings.Replace(u, fmt.Sprintf(":%s", k), v, -1)
	}

	// Add query params
	if q := p.QueryParams.Encode(); q != "" {
		sep := "?"
		if strings.Contains(u, "?") {
			sep = "&"
		}
		u = u + sep + q
	}
	return u
}

func (c *Client) createHTTPRequest(endpointURL string, body []byte) (*http.Request, error) {
	var b io.Reader
	if len(body) > 0 {
		b = strings.NewReader(string(body))
	}

	r, err := http.NewRequest(c.method, endpointURL, b)
	return r, err
}

// setHeader sets the request headers based on the resource client configuration and payload
func (c *Client) setHeader(r *http.Request, p Payload) {
	if c.contentType != "" {
		r.Header.Set("Content-Type", c.contentType)
	}

	// Set default request headers
	for k, v := range c.header.values { // Resource client default headers
		r.Header.Set(k, v)
	}
	for k, v := range p.Header { // Payload headers
		r.Header.Set(k, v)
	}
}

func (c *Client) execute(
	ctx context.Context,
	endpointURL string,
	p Payload,
) (Response, error) {
	var resultResp Response
	var attempts int

	if err := execWithRetryFunc(ctx, c.timeoutAndRetryOption.maxRetries, c.timeoutAndRetryOption.maxWaitInclRetries,
		func() error {
			// start new attempt for the http request
			attempts++

			req, err := c.createHTTPRequest(endpointURL, p.Body) // create HTTP request
			if err != nil {
				return err
			}

			c.setHeader(req, p) // set request headers

			log.Printf("requesting to %s", endpointURL)
			resp, err := c.HTTPClient.Do(req)
			if err != nil {
				if errors.Is(ctx.Err(), context.DeadlineExceeded) {
					return backoff.Permanent(ErrOverflowMaxWait) // stop retry by returning backoff.Permanent error
				}

				// evaluate if err is caused by connection timeout
				uerr, ok := err.(*url.Error)
				if !ok || !uerr.Timeout() {
					if errors.Is(err, context.Canceled) {
						return backoff.Permanent(ErrOperationContextCanceled) // stop retry by returning backoff.Permanent error
					}
					return err
				}

				// check if we need retry on timeout or not
				if !c.timeoutAndRetryOption.onTimeout {
					return backoff.Permanent(ErrTimeout) // stop retry by returning backoff.Permanent error
				}

				return ErrTimeout
			}

			if _, ok := c.timeoutAndRetryOption.onStatusCodes[resp.StatusCode]; ok {
				log.Printf("retry on status code: (%d), attempt (%d)", resp.StatusCode, attempts)
				return fmt.Errorf("retry on status code %v", resp.StatusCode)
			}

			log.Printf("end with status code: (%d), attempt (%d)", resp.StatusCode, attempts)

			// attempt to read response body
			// we need to read the response body in the same function where we cancel the retry context
			// otherwise, for big payload, the context will be cancelled while reading the body
			log.Printf("reading body: attempt (%d)", attempts)
			respBody, err := io.ReadAll(resp.Body)
			defer resp.Body.Close()
			if err != nil {
				log.Printf("[ext_http_req] err reading from body: (%+v), attempt (%d)", err, attempts)
				if errors.Is(ctx.Err(), context.DeadlineExceeded) {
					return backoff.Permanent(ErrOverflowMaxWait) // stop retry by returning backoff.Permanent error
				}

				// evaluate if err is caused by connection timeout
				if errors.Is(err, context.Canceled) {
					return backoff.Permanent(ErrOperationContextCanceled) // stop retry by returning backoff.Permanent error
				}
				if errors.Is(err, context.DeadlineExceeded) {
					return ErrTimeout
				}

				return err
			}

			resultResp.Status = resp.StatusCode
			resultResp.Body = respBody
			resultResp.Header = resp.Header

			return nil
		}); err != nil {
		switch err {
		case context.Canceled:
			return Response{}, ErrOperationContextCanceled
		case context.DeadlineExceeded:
			return Response{}, ErrOverflowMaxWait
		default:
			return Response{}, err
		}
	}

	return resultResp, nil
}

func execWithRetryFunc(ctx context.Context, maxRetries uint64, maxWait time.Duration, retryFunc func() error) error {
	b := backoff.NewExponentialBackOff()
	// 1. Wait for 2 seconds before your first retry. (for simplicity, we're just using backoff.InitialInterval to simulate)
	b.InitialInterval = 2 * time.Second
	b.RandomizationFactor = 0
	// 2. For each following retry, the increase the wait exponentially, up to 60 seconds.
	b.MaxElapsedTime = maxWait

	return backoff.Retry(
		retryFunc,
		backoff.WithContext(
			backoff.WithMaxRetries(b, maxRetries), // 3. Set a max number of retries at which point your application considers the operation failed.
			ctx,
		),
	)
}
