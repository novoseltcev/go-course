package workers

import (
	"context"
	"sync"
)

// FanIn merges channels into one.
//
// FanIn manage output channel.
// If ctx is canceled, returned channel is closed and process is stopped.
//
// Returns channel with buffer of len(chs) that merges all channels data.
func FanIn[T any](ctx context.Context, chs ...<-chan T) <-chan T {
	finalCh := make(chan T, len(chs))

	var wg sync.WaitGroup

	for _, ch := range chs {
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
