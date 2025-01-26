package server_test

import (
	"context"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/novoseltcev/go-course/internal/schemas"
	"github.com/novoseltcev/go-course/internal/server"
	"github.com/novoseltcev/go-course/internal/storages"
	"github.com/novoseltcev/go-course/mocks"
	"github.com/novoseltcev/go-course/pkg/testutils"
	pb "github.com/novoseltcev/go-course/proto/metrics"
)

var (
	delta = int64(testutils.INT)
	value = float64(testutils.FLAOT)
)

func TestGetOne_Success(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer t.Cleanup(ctrl.Finish)

	storager := mocks.NewMockMetricStorager(ctrl)
	service := server.NewGRPCMetricsServer(logrus.New(), storager)

	got := &schemas.Metric{ID: testutils.STRING, MType: schemas.Counter, Delta: nil, Value: nil}
	storager.EXPECT().GetOne(gomock.Any(), testutils.STRING, schemas.Counter).Return(got, nil)

	resp, err := service.GetOne(context.Background(),
		&pb.GetOneRequest{Id: testutils.STRING, Type: pb.Type_counter},
	)
	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t,
		&pb.Metric{Id: testutils.STRING, Type: pb.Type_counter, Value: 0, Delta: 0},
		resp.GetMetric(),
	)
}

func TestGetOne_NotFound_FailsByNotFound(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer t.Cleanup(ctrl.Finish)

	storager := mocks.NewMockMetricStorager(ctrl)
	service := server.NewGRPCMetricsServer(logrus.New(), storager)

	storager.EXPECT().GetOne(gomock.Any(), testutils.STRING, schemas.Counter).Return(nil, storages.ErrNotFound)

	resp, err := service.GetOne(context.Background(),
		&pb.GetOneRequest{Id: testutils.STRING, Type: pb.Type_counter},
	)
	require.Error(t, err)
	require.Nil(t, resp)
	assert.Equal(t, codes.NotFound, status.Code(err))
}

func TestGetOne_FailsGetOne(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer t.Cleanup(ctrl.Finish)

	storager := mocks.NewMockMetricStorager(ctrl)
	service := server.NewGRPCMetricsServer(logrus.New(), storager)

	storager.EXPECT().GetOne(gomock.Any(), testutils.STRING, schemas.Counter).Return(nil, testutils.Err)

	resp, err := service.GetOne(context.Background(),
		&pb.GetOneRequest{Id: testutils.STRING, Type: pb.Type_counter},
	)
	require.Error(t, err)
	require.Nil(t, resp)
	assert.Equal(t, codes.Internal, status.Code(err))
}

func TestGetAll_Success(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer t.Cleanup(ctrl.Finish)

	storager := mocks.NewMockMetricStorager(ctrl)
	service := server.NewGRPCMetricsServer(logrus.New(), storager)

	got := []schemas.Metric{
		{ID: testutils.STRING, MType: schemas.Counter, Delta: &delta, Value: nil},
		{ID: testutils.STRING, MType: schemas.Gauge, Delta: nil, Value: &value},
	}
	storager.EXPECT().GetAll(gomock.Any()).Return(got, nil)

	resp, err := service.GetAll(context.Background(), &pb.GetAllRequest{})
	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.ElementsMatch(t, []*pb.Metric{
		{Id: testutils.STRING, Type: pb.Type_counter, Delta: delta, Value: 0},
		{Id: testutils.STRING, Type: pb.Type_gauge, Delta: 0, Value: value},
	}, resp.GetMetrics())
}

func TestGetAll_FailsGetAll(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer t.Cleanup(ctrl.Finish)

	storager := mocks.NewMockMetricStorager(ctrl)
	service := server.NewGRPCMetricsServer(logrus.New(), storager)

	storager.EXPECT().GetAll(gomock.Any()).Return(nil, testutils.Err)

	resp, err := service.GetAll(context.Background(), &pb.GetAllRequest{})
	require.Error(t, err)
	require.Nil(t, resp)
	assert.Equal(t, codes.Internal, status.Code(err))
}

func TestUpdate_Success(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer t.Cleanup(ctrl.Finish)

	storager := mocks.NewMockMetricStorager(ctrl)
	service := server.NewGRPCMetricsServer(logrus.New(), storager)

	storager.EXPECT().Save(gomock.Any(),
		&schemas.Metric{ID: testutils.STRING, MType: schemas.Counter, Delta: &delta},
	).Return(nil).Times(1)

	_, err := service.Update(context.Background(), &pb.UpdateRequest{
		Metric: &pb.Metric{Id: testutils.STRING, Type: pb.Type_counter, Delta: delta},
	})
	assert.NoError(t, err)
}

func TestUpdate_FailsSave(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer t.Cleanup(ctrl.Finish)

	storager := mocks.NewMockMetricStorager(ctrl)
	service := server.NewGRPCMetricsServer(logrus.New(), storager)

	storager.EXPECT().Save(gomock.Any(),
		&schemas.Metric{ID: testutils.STRING, MType: schemas.Counter, Delta: &delta},
	).Return(testutils.Err).Times(1)

	_, err := service.Update(context.Background(), &pb.UpdateRequest{
		Metric: &pb.Metric{Id: testutils.STRING, Type: pb.Type_counter, Delta: delta},
	})
	assert.Equal(t, codes.Internal, status.Code(err))
}

func TestUpdateBatch(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer t.Cleanup(ctrl.Finish)

	storager := mocks.NewMockMetricStorager(ctrl)
	service := server.NewGRPCMetricsServer(logrus.New(), storager)

	storager.EXPECT().SaveBatch(gomock.Any(), []schemas.Metric{
		{ID: testutils.STRING, MType: schemas.Counter, Delta: &delta, Value: nil},
		{ID: testutils.STRING, MType: schemas.Gauge, Delta: nil, Value: &value},
	}).Return(nil).Times(1)

	_, err := service.UpdateBatch(context.Background(), &pb.UpdateBatchRequest{
		Metrics: []*pb.Metric{
			{Id: testutils.STRING, Type: pb.Type_counter, Delta: delta, Value: 0},
			{Id: testutils.STRING, Type: pb.Type_gauge, Delta: 0, Value: value},
		},
	})
	assert.NoError(t, err)
}

func TestUpdateBatch_FailsSaveBatch(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer t.Cleanup(ctrl.Finish)

	storager := mocks.NewMockMetricStorager(ctrl)
	service := server.NewGRPCMetricsServer(logrus.New(), storager)

	storager.EXPECT().SaveBatch(gomock.Any(), []schemas.Metric{
		{ID: testutils.STRING, MType: schemas.Counter, Delta: &delta, Value: nil},
		{ID: testutils.STRING, MType: schemas.Gauge, Delta: nil, Value: &value},
	}).Return(testutils.Err).Times(1)

	_, err := service.UpdateBatch(context.Background(), &pb.UpdateBatchRequest{
		Metrics: []*pb.Metric{
			{Id: testutils.STRING, Type: pb.Type_counter, Delta: delta, Value: 0},
			{Id: testutils.STRING, Type: pb.Type_gauge, Delta: 0, Value: value},
		},
	})
	assert.Equal(t, codes.Internal, status.Code(err))
}
