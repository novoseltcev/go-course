package utils

import (
	"context"
	"errors"
	"time"

	"github.com/avast/retry-go/v4"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
)

func RetryPgSelect[T any](ctx context.Context, retryableFunc retry.RetryableFuncWithData[T], retries uint) (T, error) {
	return retry.DoWithData(
		retryableFunc,
		retry.RetryIf(func(err error) bool {
			var pgErr *pgconn.PgError

			return errors.As(err, &pgErr) && pgerrcode.IsConnectionException(pgErr.Code)
		}),
		retry.Attempts(retries),
		retry.DelayType(func(n uint, _ error, _ *retry.Config) time.Duration {
			return []time.Duration{time.Second, 3 * time.Second, 5 * time.Second}[n]
		}),
		retry.Context(ctx),
	)
}

func RetryPgExec(ctx context.Context, retryableFunc retry.RetryableFunc, retries uint) error {
	return retry.Do(
		retryableFunc,
		retry.RetryIf(func(err error) bool {
			var pgErr *pgconn.PgError

			return errors.As(err, &pgErr) && pgerrcode.IsConnectionException(pgErr.Code)
		}),
		retry.Attempts(retries),
		retry.DelayType(func(n uint, _ error, _ *retry.Config) time.Duration {
			return []time.Duration{time.Second, 3 * time.Second, 5 * time.Second}[n]
		}),
		retry.Context(ctx),
	)
}
