package workers_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/novoseltcev/go-course/pkg/workers"
)

func TestFanIn(t *testing.T) {
	t.Parallel()

	ch1 := make(chan int)
	ch2 := make(chan int)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)

	defer cancel()

	resCh := workers.FanIn[int](ctx, ch1, ch2)
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

	go func() {
		val1Cnt := 0
		val2Cnt := 0

		for range cnt {
			for range 4 {
				tmp := <-resCh
				if tmp == testValue1 {
					val1Cnt++
				}

				if tmp == testValue2 {
					val2Cnt++
				}
			}
		}

		assert.Equal(t, 2*cnt, val1Cnt)
		assert.Equal(t, 2*cnt, val2Cnt)
	}()

	<-ctx.Done()
}
