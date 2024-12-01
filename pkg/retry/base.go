// Package retry presents decorators for retrying operations.
package retry

import (
	"context"
	"time"

	"github.com/avast/retry-go/v4"
)

// Do is wrapper for retry.Do function with.
//
// Configured by passed Options.
// Returns last error if the all attempts failed.
func Do(
	ctx context.Context,
	retryableFunc retry.RetryableFunc,
	opts *Options,
) error {
	if opts == nil {
		opts = &Options{}
	}

	return retry.Do(
		retryableFunc,
		retry.Attempts(opts.TotalAttempts()),
		retry.DelayType(func(n uint, _ error, _ *retry.Config) time.Duration {
			return opts.GetAttemptDelay(n)
		}),
		retry.Context(ctx),
	)
}

// DoWithData retries the function with the given options.
//
// The function must return a value and an error.
// Returns last error if the all attempts failed.
func DoWithData[T any](
	ctx context.Context,
	retryableFunc retry.RetryableFuncWithData[T],
	opts *Options,
) (T, error) {
	if opts == nil {
		opts = &Options{}
	}

	return retry.DoWithData(
		retryableFunc,
		retry.Attempts(opts.TotalAttempts()),
		retry.DelayType(func(n uint, _ error, _ *retry.Config) time.Duration {
			return opts.GetAttemptDelay(n)
		}),
		retry.Context(ctx),
	)
}
