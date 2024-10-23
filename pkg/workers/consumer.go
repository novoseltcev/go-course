// Template workers with go-channels, implement some patterns.
package workers

import (
	"context"
	"time"

	log "github.com/sirupsen/logrus"
)

// Consumer is consumer that processes data from channel jobsCh and calls consumeFunc for each value.
//
// May be used as goroutine or as a function.
//
// If ctx is canceled, consumer stops.
// If consumeFunc returns error, consumer skips value and continues.
//
// Usage:
//
//	go Consumer(ctx, jobsCh, func(ctx context.Context, value T) error {
//	    // do something with value
//	})
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

// AntiFraudConsumer is a consumer that processes tasks in batches.
//
// May be used as goroutine or as a function.
// Use rateLimit to set the maximum time between processing batches.
//
// If ctx is canceled, consumer stops.
// If consumeFunc returns error, consumer don't clear the buffer and try to process this batch again.
//
// Usage:
//
//	go AntiFraudConsumer(ctx, jobsCh, func(ctx context.Context, buf []T) error {
//	    // do something with buf
//	}, rateLimit)
func AntiFraudConsumer[T any](
	ctx context.Context,
	jobsCh <-chan T,
	consumeFunc func(ctx context.Context, buf []T) error,
	rateLimit time.Duration,
) {
	buf := []T{}             // buffer for tasks
	lstConsume := time.Now() // last time when the buffer was consumed

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
