// nolint: paralleltest
package workers_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/novoseltcev/go-course/pkg/workers"
)

func TestConsumerSuccess(t *testing.T) {
	ch := make(chan int)
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)

	defer cancel()

	calls := 0
	sent := 0

	go workers.Consumer(ctx, ch, func(_ context.Context, _ int) error {
		calls++

		return nil
	})

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case ch <- 1:
				sent++
			}
		}
	}()

	<-ctx.Done()
	assert.InDelta(t, sent, calls, 1.)
}

func TestConsumerErr(t *testing.T) {
	t.Parallel()

	ch := make(chan int)
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)

	defer cancel()

	go workers.Consumer(ctx, ch, func(_ context.Context, _ int) error {
		return errSome
	})

	ch <- 1

	<-ctx.Done()
}

func TestAntiFraudConsumerOnce(t *testing.T) {
	ch := make(chan int)
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond+500*time.Microsecond)

	defer cancel()

	calls := 0

	go workers.AntiFraudConsumer(ctx, ch, func(_ context.Context, _ []int) error {
		calls++

		return nil
	}, time.Millisecond)

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case ch <- 1:
			}
		}
	}()

	<-ctx.Done()
	assert.Equal(t, 1, calls)
}

func TestAntiFraudConsumerMany(t *testing.T) {
	ch := make(chan int)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Millisecond+500*time.Microsecond)

	defer cancel()

	calls := 0

	go workers.AntiFraudConsumer(ctx, ch, func(_ context.Context, _ []int) error {
		calls++

		return nil
	}, time.Millisecond)

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case ch <- 1:
			}
		}
	}()

	<-ctx.Done()
	assert.Equal(t, 2, calls)
}

func TestAntiFraudConsumerNotCalledFn(t *testing.T) {
	ch := make(chan int)
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)

	defer cancel()

	calls := 0

	go workers.AntiFraudConsumer(ctx, ch, func(_ context.Context, _ []int) error {
		calls++

		return nil
	}, 2*time.Millisecond)

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case ch <- 1:
			}
		}
	}()

	<-ctx.Done()
	assert.Equal(t, 0, calls)
}

func TestAntiFraudConsumerErr(t *testing.T) {
	ch := make(chan int)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Millisecond)

	defer cancel()

	calls := 0

	go workers.AntiFraudConsumer(ctx, ch, func(_ context.Context, _ []int) error {
		calls++

		return errSome
	}, time.Millisecond)

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case ch <- 1:
			}
		}
	}()

	<-ctx.Done()
	assert.GreaterOrEqual(t, calls, 2)
}
