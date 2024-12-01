package schemas_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/novoseltcev/go-course/internal/schemas"
)

func TestMetric_Validate(t *testing.T) {
	t.Parallel()

	var (
		testID    = "test"
		testDelta = int64(1)
		testValue = float64(1)
	)

	tests := []struct {
		name    string
		metric  *schemas.Metric
		wantErr error
	}{
		{
			name: "valid gauge",
			metric: &schemas.Metric{
				ID:    testID,
				MType: schemas.Gauge,
				Value: &testValue,
			},
			wantErr: nil,
		},
		{
			name: "valid counter",
			metric: &schemas.Metric{
				ID:    testID,
				MType: schemas.Counter,
				Delta: &testDelta,
			},
			wantErr: nil,
		},
		{
			name:    "empty id",
			metric:  &schemas.Metric{},
			wantErr: errors.New("id is required"),
		},
		{
			name: "invalid type",
			metric: &schemas.Metric{
				ID:    testID,
				MType: "unknown",
			},
			wantErr: errors.New("type is invalid"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tt.wantErr, tt.metric.Validate())
		})
	}
}
