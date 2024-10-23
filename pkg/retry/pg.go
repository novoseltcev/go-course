package retry

import (
	"context"
	"errors"
	"time"

	"github.com/avast/retry-go/v4"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
)

// PgExec look like retry.Do but specified for postgresql.
//
// Retry only if error is a connection exception.
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

// PgSelect look like retry.DoWithData but specified for postgresql.
//
// Retry only if error is a connection exception.
func PgSelect[T any](
	ctx context.Context,
	retryableFunc retry.RetryableFuncWithData[T],
	opts *Options,
) (T, error) {
	if opts == nil {
		opts = &Options{}
	}

	return retry.DoWithData(
		retryableFunc,
		retry.RetryIf(func(err error) bool {
			var pgErr *pgconn.PgError

			return errors.As(err, &pgErr) && pgerrcode.IsConnectionException(pgErr.Code)
		}),
		retry.Attempts(opts.TotalAttempts()),
		retry.DelayType(func(n uint, _ error, _ *retry.Config) time.Duration {
			return opts.GetAttemptDelay(n)
		}),
		retry.Context(ctx),
	)
}
