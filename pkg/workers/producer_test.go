package workers_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/novoseltcev/go-course/pkg/workers"
)

func TestProducerSuccess(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)

	defer cancel()

	testData := []int{1, 2, 3}
	produced := make([]int, 0)

	ch := workers.Producer(ctx, func(_ context.Context) ([]int, error) {
		return testData, nil
	}, time.Second)

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

	require.NotEmpty(t, produced)
	assert.InDeltaSlice(t, testData, produced, 1)
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
