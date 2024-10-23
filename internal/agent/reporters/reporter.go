package reporters

import (
	"context"

	"github.com/novoseltcev/go-course/internal/schemas"
)

type Reporter interface {
	Report(ctx context.Context, metrics []schemas.Metric) error
}
