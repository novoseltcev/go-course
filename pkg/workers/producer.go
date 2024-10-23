package workers

import (
	"context"
	"time"
)

// Producer returns channel that produces data from dataFunc with delay.
//
// Producer manage output channel.
// If ctx is canceled, returned channel is closed.
// If dataFunc returns error, produced data is skipped.
func Producer[T any](
	ctx context.Context,
	dataFunc func(ctx context.Context) ([]T, error),
	delay time.Duration,
) <-chan T {
	ch := make(chan T)
	go func() {
		defer close(ch)

		for {
			values, err := dataFunc(ctx)
			if err != nil {
				continue
			}

			for _, v := range values {
				ch <- v
			}

			select {
			case <-ctx.Done():
				return
			default:
				time.Sleep(delay)
			}
		}
	}()

	return ch
}
