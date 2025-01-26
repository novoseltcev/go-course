package workers

import (
	"context"
	"time"
)

func Every(ctx context.Context, fn func(ctx context.Context) error, interval time.Duration) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			if err := fn(ctx); err != nil {
				return err
			}

			time.Sleep(interval)
		}
	}
}
