package retry

import (
	"context"
	"errors"
	"time"

	"github.com/avast/retry-go/v4"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
)

func PgSelect[T any](
	ctx context.Context,
	retryableFunc retry.RetryableFuncWithData[T],
	options *Options,
) (T, error) {
	if options == nil {
		options = &Options{}
	}

	return retry.DoWithData(
		retryableFunc,
		retry.RetryIf(func(err error) bool {
			var pgErr *pgconn.PgError

			return errors.As(err, &pgErr) && pgerrcode.IsConnectionException(pgErr.Code)
		}),
		retry.Attempts(options.TotalAttempts()),
		retry.DelayType(func(n uint, _ error, _ *retry.Config) time.Duration {
			return options.GetAttemptDelay(n)
		}),
		retry.Context(ctx),
	)
}

func PgExec(ctx context.Context, retryableFunc retry.RetryableFunc, options *Options) error {
	if options == nil {
		options = &Options{}
	}

	return retry.Do(
		retryableFunc,
		retry.RetryIf(func(err error) bool {
			var pgErr *pgconn.PgError

			return errors.As(err, &pgErr) && pgerrcode.IsConnectionException(pgErr.Code)
		}),
		retry.Attempts(options.TotalAttempts()),
		retry.DelayType(func(n uint, _ error, _ *retry.Config) time.Duration {
			return options.GetAttemptDelay(n)
		}),
		retry.Context(ctx),
	)
}
