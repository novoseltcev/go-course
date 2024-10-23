package retry

import (
	"context"
	"time"

	"github.com/avast/retry-go/v4"
)

func Do(
	ctx context.Context,
	retryableFunc retry.RetryableFunc,
	options *Options,
) error {
	if options == nil {
		options = &Options{}
	}

	return retry.Do(
		retryableFunc,
		retry.Attempts(options.TotalAttempts()),
		retry.DelayType(func(n uint, _ error, _ *retry.Config) time.Duration {
			return options.GetAttemptDelay(n)
		}),
		retry.Context(ctx),
	)
}

func DoWithData[T any](
	ctx context.Context,
	retryableFunc retry.RetryableFuncWithData[T],
	options *Options,
) (T, error) {
	if options == nil {
		options = &Options{}
	}

	return retry.DoWithData(
		retryableFunc,
		retry.Attempts(options.TotalAttempts()),
		retry.DelayType(func(n uint, _ error, _ *retry.Config) time.Duration {
			return options.GetAttemptDelay(n)
		}),
		retry.Context(ctx),
	)
}
