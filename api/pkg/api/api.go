package api

import (
	"context"
	"net/http"
	"net/url"
)

// Payload is the request payload struct representation
type Payload struct {
	// Request body
	Body []byte
	// QueryParams contains the request/query parameters
	QueryParams url.Values
	// PathVars contains the path variables used to replace placeholders
	// wrapped with {} in Client.URL
	PathVars map[string]string
	// Header contains custom request headers that will be added to the request
	// on http call.
	// The values on this field will override Client.Headers.Values
	Header map[string]string
}

// Response is the result of the http call
type Response struct {
	Status int
	Body   []byte
	Header http.Header
}

func (c *Client) Send(ctx context.Context, p Payload) (Response, error) {
	// Endpoint URL
	endpointURL := c.constructURL(p)

	// Create context with max timeout
	ctxTimeout, cancel := context.WithTimeout(ctx, c.timeoutAndRetryOption.maxWaitInclRetries)
	defer cancel()

	// HTTP operation
	resp, err := c.execute(ctxTimeout, endpointURL, p)
	if err != nil {
		return Response{}, err
	}

	return resp, nil
}
