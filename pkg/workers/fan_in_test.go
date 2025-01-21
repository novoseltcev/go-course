package workers_test

import (
	"context"
	"fmt"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/novoseltcev/go-course/pkg/testutils/helpers"
	"github.com/novoseltcev/go-course/pkg/workers"
)

func TestFanIn(t *testing.T) {
	t.Parallel()

	ch1 := make(chan int, 2)
	ch2 := make(chan int, 2)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	resCh := workers.FanIn(ctx, ch1, ch2)
	testValue1 := 1
	testValue2 := 2

	cnt := 100

	go func() {
		for range cnt {
			ch1 <- testValue1
			ch2 <- testValue1
			ch1 <- testValue2
			ch2 <- testValue2
		}
	}()

	var val1Cnt atomic.Int32
	var val2Cnt atomic.Int32

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				tmp := <-resCh
				if tmp == testValue1 {
					val1Cnt.Add(1)
				}

				if tmp == testValue2 {
					val2Cnt.Add(1)
				}
			}
		}
	}()

	<-ctx.Done()

	assert.Equal(t, 2*cnt, int(val1Cnt.Load()))
	assert.Equal(t, 2*cnt, int(val2Cnt.Load()))
}

func TestFanInWithCancelRecive(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	ch := workers.FanIn(ctx, make(chan int))

	cancel()
	helpers.AssertChannelClosed(t, ch)
}

func TestFanInWithCancelSend(t *testing.T) {
	t.Parallel()

	ch1 := make(chan int, 1)
	defer close(ch1)

	ctx, cancel := context.WithCancel(context.Background())
	ch := workers.FanIn(ctx, ch1)

	ch1 <- 1
	ch1 <- 1
	ch1 <- 1
	cancel()

	<-ch
	<-ch
	helpers.AssertChannelClosed(t, ch)
}

func TestFanInWithCloseChannels(t *testing.T) {
	t.Parallel()

	ch1 := make(chan int)
	ch2 := make(chan int)

	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)
	defer cancel()

	ch := workers.FanIn(ctx, ch1, ch2)
	ch1 <- 1
	_, ok := <-ch
	assert.True(t, ok)

	close(ch1)
	ch2 <- 1
	_, ok = <-ch
	assert.True(t, ok)

	close(ch2)
	helpers.AssertChannelClosed(t, ch)
}

func ExampleFanIn() {
	ch1 := make(chan int)
	ch2 := make(chan int)
	defer close(ch1)
	defer close(ch2)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Millisecond)
	defer cancel()

	ch := workers.FanIn(ctx, ch1, ch2)

	ch1 <- 1
	ch2 <- 2

	for {
		select {
		case <-ctx.Done():
			return
		case data, ok := <-ch:
			if ok {
				fmt.Println(data)
			}
		}
	}

	// Output:
	// 1
	// 2
}
