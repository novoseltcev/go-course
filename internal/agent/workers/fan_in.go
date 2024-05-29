package workers

import (
	"context"
	"sync"

	"github.com/novoseltcev/go-course/internal/model"
)

func FanIn(ctx context.Context, resultChs ...<-chan model.Metric) <-chan model.Metric {
    finalCh := make(chan model.Metric, len(resultChs))
    var wg sync.WaitGroup

    for _, ch := range resultChs {
        chClosure := ch
        wg.Add(1)
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
