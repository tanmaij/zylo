package api

import (
	"context"
	"errors"
	"fmt"
)

var (
	// ErrTimeout means timeout
	ErrTimeout = errors.New("timeout")
	// ErrOverflowMaxWait means overflow the max wait limit
	ErrOverflowMaxWait = errors.New("overflow max wait")
	// ErrOperationContextCanceled means the given context for the operation was canceled
	ErrOperationContextCanceled = fmt.Errorf("operation context canceled: %w", context.Canceled)
	// ErrTimeoutAndRetryOptionInvalid means the retry config is invalid
	ErrTimeoutAndRetryOptionInvalid = errors.New("retry config invalid")
)
