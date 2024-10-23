package workers

import (
	"context"
	"time"

	log "github.com/sirupsen/logrus"
)

func Consumer[T any](
	ctx context.Context,
	jobsCh <-chan T,
	consumeFunc func(ctx context.Context, value T) error,
) {
	for v := range jobsCh {
		select {
		case <-ctx.Done():
			return
		default:
			if err := consumeFunc(ctx, v); err != nil {
				log.WithError(err).Error("crash worker")
			}
		}
	}
}

func AntiFraudConsumer[T any](
	ctx context.Context,
	jobsCh <-chan T,
	consumeFunc func(ctx context.Context, buf []T) error,
	rateLimit time.Duration,
) {
	buf := []T{}
	lstConsume := time.Now()

	Consumer(ctx, jobsCh, func(ctx context.Context, v T) error {
		buf = append(buf, v)

		if time.Since(lstConsume) > rateLimit {
			if err := consumeFunc(ctx, buf); err != nil {
				return err
			}

			buf = make([]T, 0)
			lstConsume = time.Now()
		}

		return nil
	})
}
