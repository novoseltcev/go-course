package workers

import (
	"context"
	"time"

	log "github.com/sirupsen/logrus"
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
			select {
			case <-ctx.Done():
				return
			default:
				time.Sleep(delay)
			}

			values, err := dataFunc(ctx)
			if err != nil {
				log.WithError(err).WithField("func", dataFunc).Error("failed to get data")

				continue
			}

			for _, v := range values {
				ch <- v
			}
		}
	}()

	return ch
}
