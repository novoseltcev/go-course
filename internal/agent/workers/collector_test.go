package workers_test

import (
	"context"
	"testing"
	"time"

	"github.com/novoseltcev/go-course/internal/agent/workers"
)

func TestCollectMetrics(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)

	defer cancel()

	ch := workers.CollectMetrics(ctx, time.Second)

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case metric := <-ch:
				t.Log(metric)
			}
		}
	}()

	<-ctx.Done()
}

func TestCollectCoreMetrics(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)

	defer cancel()

	ch := workers.CollectCoreMetrics(ctx, time.Second)

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case metric := <-ch:
				t.Log(metric)
			}
		}
	}()

	<-ctx.Done()
}
