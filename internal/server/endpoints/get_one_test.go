package endpoints_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	gomock "go.uber.org/mock/gomock"

	"github.com/novoseltcev/go-course/internal/schemas"
	"github.com/novoseltcev/go-course/internal/server/endpoints"
	"github.com/novoseltcev/go-course/internal/storages"
	"github.com/novoseltcev/go-course/mocks"
)

// nolint: gochecknoglobals
var (
	errTest   = errors.New("test error")
	testID    = "test"
	testValue = 10.123
	testDelta = int64(10)
	testGauge = schemas.Metric{
		ID:    testID,
		MType: schemas.Gauge,
		Value: &testValue,
	}
	testCounter = schemas.Metric{
		ID:    testID,
		MType: schemas.Counter,
		Delta: &testDelta,
	}
)

func TestGetOneMetric(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	storager := mocks.NewMockMetricStorager(ctrl)
	router := endpoints.NewAPIRouter(storager)

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
			got{id: testID, Type: schemas.Gauge},
			&want{&testGauge, nil},
			http.StatusOK,
			`10.123`,
		},
		{
			"success counter",
			got{id: testID, Type: schemas.Counter},
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
			got{id: testID, Type: schemas.Gauge},
			&want{nil, errTest},
			http.StatusInternalServerError,
			"failed to get metric\n",
		},
		{
			name: "invalid metric type",
			got:  got{id: testID, Type: "unknown"},
			code: http.StatusBadRequest,
			body: "type is invalid\n",
		},
		{
			"failed serialization unknown",
			got{id: testID, Type: schemas.Counter},
			&want{&schemas.Metric{ID: testID, MType: "unknown"}, nil},
			http.StatusInternalServerError,
			"failed to serialize\n",
		},
		{
			"failed serialization empty delta",
			got{id: testID, Type: schemas.Counter},
			&want{&schemas.Metric{ID: testID, MType: schemas.Counter}, nil},
			http.StatusInternalServerError,
			"failed to serialize\n",
		},
		{
			"failed serialization empty value",
			got{id: testID, Type: schemas.Gauge},
			&want{&schemas.Metric{ID: testID, MType: schemas.Gauge}, nil},
			http.StatusInternalServerError,
			"failed to serialize\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if tt.want != nil {
				storager.EXPECT().GetOne(gomock.Any(), tt.got.id, tt.got.Type).Return(tt.want.metric, tt.want.err)
			}

			req := httptest.NewRequest(http.MethodGet, "/value/"+tt.got.Type+"/"+tt.got.id, http.NoBody)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.code, w.Code)
			assert.Equal(t, tt.body, w.Body.String())
		})
	}
}
