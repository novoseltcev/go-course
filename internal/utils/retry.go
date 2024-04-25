package utils

import (
	"context"
	"errors"
	"time"

	"github.com/avast/retry-go/v4"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
)

func RetryPgSelect[T any](ctx context.Context, retryableFunc retry.RetryableFuncWithData[T]) (T, error) {
	return retry.DoWithData(
		retryableFunc,
		retry.RetryIf(func(err error) bool {
			var pgErr *pgconn.PgError
			return errors.As(err, &pgErr) && pgerrcode.IsConnectionException(pgErr.Code)
		}),
		retry.Attempts(3),
		retry.DelayType(func(n uint, err error, config *retry.Config) time.Duration {
			return []time.Duration{time.Second, 3 * time.Second, 5 * time.Second}[n]
		}),
		retry.Context(ctx),
	)
}

func RetryPgExec(ctx context.Context, retryableFunc retry.RetryableFunc) (error) {
	return retry.Do(
		retryableFunc,
		retry.RetryIf(func(err error) bool {
			var pgErr *pgconn.PgError
			return errors.As(err, &pgErr) && pgerrcode.IsConnectionException(pgErr.Code)
		}),
		retry.Attempts(3),
		retry.DelayType(func(n uint, err error, config *retry.Config) time.Duration {
			return []time.Duration{time.Second, 3 * time.Second, 5 * time.Second}[n]
		}),
		retry.Context(ctx),
	)
}
