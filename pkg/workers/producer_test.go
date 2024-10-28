package workers_test

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/novoseltcev/go-course/pkg/workers"
)

func ExampleProducer() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)

	defer cancel()

	ch := workers.Producer(ctx, func(_ context.Context) ([]int, error) {
		return []int{1, 2, 3}, nil
	}, time.Second)

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case v := <-ch:
				fmt.Println(v)
			}
		}
	}()

	<-ctx.Done()

	// Output:
	// 1
	// 2
	// 3
}

var errSome = errors.New("some error")

func TestProducerErr(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)

	defer cancel()

	ch := workers.Producer(ctx, func(_ context.Context) ([]int, error) {
		return nil, errSome
	}, time.Second)

	produced := make([]int, 0)

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case v := <-ch:
				produced = append(produced, v)
			}
		}
	}()

	<-ctx.Done()

	assert.Empty(t, produced)
}
