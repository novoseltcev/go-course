package endpoints_test

import (
	"net/http"
	"testing"

	"github.com/steinfletcher/apitest"
	gomock "go.uber.org/mock/gomock"

	"github.com/novoseltcev/go-course/internal/schemas"
	"github.com/novoseltcev/go-course/internal/server/endpoints"
	"github.com/novoseltcev/go-course/internal/storages"
	"github.com/novoseltcev/go-course/mocks"
	"github.com/novoseltcev/go-course/pkg/testutils"
)

var (
	testValue = 10.123
	testDelta = int64(10)
	testGauge = schemas.Metric{
		ID:    testutils.STRING,
		MType: schemas.Gauge,
		Value: &testValue,
	}
	testCounter = schemas.Metric{
		ID:    testutils.STRING,
		MType: schemas.Counter,
		Delta: &testDelta,
	}
)

func TestGetOneMetric(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	type got struct {
		id   string
		Type string
	}

	type want struct {
		metric *schemas.Metric
		err    error
	}

	tests := []struct {
		name string
		got
		want *want
		code int
		body string
	}{
		{
			"success gauge",
			got{id: testutils.STRING, Type: schemas.Gauge},
			&want{&testGauge, nil},
			http.StatusOK,
			`10.123`,
		},
		{
			"success counter",
			got{id: testutils.STRING, Type: schemas.Counter},
			&want{&testCounter, nil},
			http.StatusOK,
			`10`,
		},
		{
			"metric not found",
			got{id: "unknown", Type: schemas.Gauge},
			&want{nil, storages.ErrNotFound},
			http.StatusNotFound,
			"metric not found\n",
		},
		{
			"failed get",
			got{id: testutils.STRING, Type: schemas.Gauge},
			&want{nil, testutils.Err},
			http.StatusInternalServerError,
			"failed to get metric\n",
		},
		{
			name: "invalid metric type",
			got:  got{id: testutils.STRING, Type: "unknown"},
			code: http.StatusBadRequest,
			body: "metric-id validator: type is invalid\n",
		},
		{
			"failed serialization unknown",
			got{id: testutils.STRING, Type: schemas.Counter},
			&want{&schemas.Metric{ID: testutils.STRING, MType: "unknown"}, nil},
			http.StatusInternalServerError,
			"failed to serialize\n",
		},
		{
			"failed serialization empty delta",
			got{id: testutils.STRING, Type: schemas.Counter},
			&want{&schemas.Metric{ID: testutils.STRING, MType: schemas.Counter}, nil},
			http.StatusInternalServerError,
			"failed to serialize\n",
		},
		{
			"failed serialization empty value",
			got{id: testutils.STRING, Type: schemas.Gauge},
			&want{&schemas.Metric{ID: testutils.STRING, MType: schemas.Gauge}, nil},
			http.StatusInternalServerError,
			"failed to serialize\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			storager := mocks.NewMockMetricStorager(ctrl)
			if tt.want != nil {
				storager.EXPECT().GetOne(gomock.Any(), tt.got.id, tt.got.Type).Return(tt.want.metric, tt.want.err)
			}

			apitest.New().
				Handler(endpoints.NewAPIRouter(storager)).
				Getf("/value/%s/%s", tt.got.Type, tt.got.id).
				Expect(t).
				Status(tt.code).
				Body(tt.body).
				End()
		})
	}
}
