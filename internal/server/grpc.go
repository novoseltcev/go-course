package server

import (
	"context"
	"errors"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/novoseltcev/go-course/internal/schemas"
	"github.com/novoseltcev/go-course/internal/storages"
	pb "github.com/novoseltcev/go-course/proto/metrics"
)

type MetricsServer struct {
	pb.UnimplementedMetricsServiceServer

	logger   *logrus.Logger
	storager storages.MetricStorager
}

func NewGRPCMetricsServer(logger *logrus.Logger, storager storages.MetricStorager) *MetricsServer {
	return &MetricsServer{logger: logger, storager: storager}
}

func (s *MetricsServer) GetOne(ctx context.Context, r *pb.GetOneRequest) (*pb.GetOneResponse, error) {
	var resp pb.GetOneResponse

	metric, err := s.storager.GetOne(ctx, r.GetId(), r.GetType().String())
	if err != nil {
		if errors.Is(err, storages.ErrNotFound) {
			return nil, status.Error(codes.NotFound, err.Error())
		}

		s.logger.WithError(err).Error("failed to get metric")

		return nil, status.Error(codes.Internal, "failed to get metric")
	}

	resp.Metric, err = mapSchemaToMsg(metric)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to map schema to message")
	}

	return &resp, nil
}

func (s *MetricsServer) GetAll(ctx context.Context, _ *pb.GetAllRequest) (*pb.GetAllResponse, error) {
	var resp pb.GetAllResponse

	metrics, err := s.storager.GetAll(ctx)
	if err != nil {
		s.logger.WithError(err).Error("failed to get metrics")

		return nil, status.Error(codes.Internal, "failed to get metrics")
	}

	resp.Metrics = make([]*pb.Metric, len(metrics))
	for i, metric := range metrics {
		resp.Metrics[i], err = mapSchemaToMsg(&metric)
		if err != nil {
			return nil, status.Error(codes.Internal, "failed to map schema to message")
		}
	}

	return &resp, nil
}

func (s *MetricsServer) Update(ctx context.Context, r *pb.UpdateRequest) (*pb.UpdateResponse, error) {
	var response pb.UpdateResponse

	schema, err := mapMsgToSchema(r.GetMetric())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "failed to map message to schema")
	}

	if err := s.storager.Save(ctx, schema); err != nil {
		s.logger.WithError(err).Error("failed to save metric")

		return nil, status.Error(codes.Internal, "failed to save metric")
	}

	return &response, nil
}

func (s *MetricsServer) UpdateBatch(ctx context.Context, req *pb.UpdateBatchRequest) (*pb.UpdateBatchResponse, error) {
	var response pb.UpdateBatchResponse

	metrics := make([]schemas.Metric, len(req.GetMetrics()))

	for i, metric := range req.GetMetrics() {
		schema, err := mapMsgToSchema(metric)
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, "failed to map message to schema")
		}

		metrics[i] = *schema
	}

	if err := s.storager.SaveBatch(ctx, metrics); err != nil {
		s.logger.WithError(err).Error("failed to save metrics")

		return nil, status.Error(codes.Internal, "failed to save metrics")
	}

	return &response, nil
}

func mapMsgToSchema(msg *pb.Metric) (*schemas.Metric, error) {
	result := schemas.Metric{
		ID:    msg.GetId(),
		MType: msg.GetType().String(),
	}

	if msg.GetType() == pb.Type_counter {
		result.Delta = &msg.Delta
	} else if msg.GetType() == pb.Type_gauge {
		result.Value = &msg.Value
	}

	return &result, result.Validate()
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
