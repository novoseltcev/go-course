package reporters

import (
	"context"

	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/novoseltcev/go-course/internal/schemas"
	pb "github.com/novoseltcev/go-course/proto/metrics"
)

//go:generate mockgen -source=../../../proto/metrics/metrics_grpc.pb.go -destination=./grpc_mock_test.go -package=reporters_test -typed

type GRPCReporter struct {
	client pb.MetricsServiceClient
}

func NewGRPCReporter(client pb.MetricsServiceClient) *GRPCReporter {
	return &GRPCReporter{client: client}
}

func (rc *GRPCReporter) Report(ctx context.Context, metrics []schemas.Metric) error {
	if len(metrics) == 0 {
		return nil
	}

	var req pb.UpdateBatchRequest
	req.Metrics = make([]*pb.Metric, len(metrics))

	for i, metric := range metrics {
		var err error

		req.Metrics[i], err = mapSchemaToMsg(&metric)
		if err != nil {
			return err
		}
	}

	_, err := rc.client.UpdateBatch(ctx, &req)
	if err != nil {
		switch status.Code(err) { // nolint: exhaustive
		case codes.Unavailable:
			log.WithError(err).Error("failed to send metrics")
		case codes.Internal:
			log.WithError(err).Error("internal server error")
		case codes.InvalidArgument:
			log.WithError(err).Error("invalid argument")
		default:
			log.WithError(err).Error("unknown error")
		}

		return err
	}

	log.Info("report successfully sent")

	return nil
}

func mapSchemaToMsg(metric *schemas.Metric) (*pb.Metric, error) {
	Type, ok := pb.Type_value[metric.MType]
	if !ok {
		return nil, schemas.ErrInvalidType
	}

	result := pb.Metric{
		Id:   metric.ID,
		Type: pb.Type(Type),
	}

	if metric.Value != nil {
		result.Value = *metric.Value
	}

	if metric.Delta != nil {
		result.Delta = *metric.Delta
	}

	return &result, nil
}
