package workers

import (
	"context"
	"sync"

	"github.com/novoseltcev/go-course/internal/schema"
)

func FanIn(ctx context.Context, resultChs ...<-chan schema.Metric) <-chan schema.Metric {
	finalCh := make(chan schema.Metric, len(resultChs))

	var wg sync.WaitGroup

	for _, ch := range resultChs {
		wg.Add(1)

		chClosure := ch

		go func() {
			defer wg.Done()

			for data := range chClosure {
				select {
				case <-ctx.Done():
					return
				case finalCh <- data:
				}
			}
		}()
	}

	go func() {
		wg.Wait()
		close(finalCh)
	}()

	return finalCh
}
