package workers

import (
	"context"
	"testing"
	"time"
)

func TestCollectMetrics(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second * 1)
	defer cancel()

	ch := CollectMetrics(ctx, time.Second)

	for {
		select {
		case <-ctx.Done():
			return
		case metric := <-ch:
			t.Log(metric)
		}
	}
}

func TestCollectCoreMetrics(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second * 1)
	defer cancel()

	ch := CollectCoreMetrics(ctx, time.Second)

	for {
		select {
		case <-ctx.Done():
			return
		case metric := <-ch:
			t.Log(metric)
		}
	}

}
