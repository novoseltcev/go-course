package workers

import (
	"context"
	"time"
)

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
