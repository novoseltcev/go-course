// nolint: paralleltest
package reporters_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/novoseltcev/go-course/internal/agent/reporters"
	"github.com/novoseltcev/go-course/internal/schemas"
	"github.com/novoseltcev/go-course/pkg/testutils"
	pb "github.com/novoseltcev/go-course/proto/metrics"
)

func TestGRPCReporter_Report_Empty_Success(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	reporter := reporters.NewGRPCReporter(NewMockMetricsServiceClient(ctrl))
	assert.NoError(t, reporter.Report(context.TODO(), []schemas.Metric{}))
}

func TestGRPCReporter_Report_Success(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	tests := []struct {
		name string
		got  []schemas.Metric
		want []*pb.Metric
	}{
		{
			name: "gauge",
			got: []schemas.Metric{
				{
					ID:    testutils.STRING,
					MType: schemas.Gauge,
					Value: &value,
				},
			},
			want: []*pb.Metric{
				{
					Id:    testutils.STRING,
					Type:  pb.Type_gauge,
					Value: value,
				},
			},
		},
		{
			name: "counter",
			got: []schemas.Metric{
				{
					ID:    testutils.STRING,
					MType: schemas.Counter,
					Delta: &delta,
				},
			},
			want: []*pb.Metric{
				{
					Id:    testutils.STRING,
					Type:  pb.Type_counter,
					Delta: delta,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			client := NewMockMetricsServiceClient(ctrl)
			reporter := reporters.NewGRPCReporter(client)

			client.EXPECT().
				UpdateBatch(gomock.Any(), gomock.Any()).
				Return(&pb.UpdateBatchResponse{}, nil).
				Times(1)

			assert.NoError(t, reporter.Report(context.TODO(), tt.got))
		})
	}
}

func TestGRPCReporter_Report_Fails(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	tests := []struct {
		name string
		code codes.Code
	}{
		{
			name: "unavailable",
			code: codes.Unavailable,
		},
		{
			name: "internal",
			code: codes.Internal,
		},
		{
			name: "invalid argument",
			code: codes.InvalidArgument,
		},
		{
			name: "unknown",
			code: codes.Unknown,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			client := NewMockMetricsServiceClient(ctrl)
			reporter := reporters.NewGRPCReporter(client)

			client.EXPECT().
				UpdateBatch(gomock.Any(), gomock.Any()).
				Return(nil, status.Error(tt.code, testutils.Err.Error())).
				Times(1)

			err := reporter.Report(context.TODO(), testMetrics)
			assert.Equal(t, tt.code, status.Code(err))
		})
	}
}
