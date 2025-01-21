package workers_test

import (
	"context"
	"fmt"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"golang.org/x/sync/errgroup"

	"github.com/novoseltcev/go-course/pkg/testutils"
	"github.com/novoseltcev/go-course/pkg/workers"
)

func ExampleEvery() {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Millisecond)
	defer cancel()

	var i atomic.Int64
	go workers.Every(ctx, func(_ context.Context) error {
		fmt.Println(i.Add(1))

		return nil
	}, time.Millisecond)

	<-ctx.Done()

	// Output:
	// 1
	// 2
	// 3
}

func TestEveryErr(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)
	defer cancel()

	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		return workers.Every(ctx, func(_ context.Context) error { return testutils.Err }, time.Microsecond)
	})

	err := g.Wait()
	assert.ErrorIs(t, err, testutils.Err)
}
