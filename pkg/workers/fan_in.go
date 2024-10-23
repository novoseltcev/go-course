package workers

import (
	"context"
	"sync"
)

func FanIn[T any](ctx context.Context, resultChs ...<-chan T) <-chan T {
	finalCh := make(chan T, len(resultChs))

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
