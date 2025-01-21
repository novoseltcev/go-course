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

func TestGetOneMetricFromJSON(t *testing.T) {
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
			`{"value":10.123,"id":"string","type":"gauge"}`,
		},
		{
			"success counter",
			got{id: testutils.STRING, Type: schemas.Counter},
			&want{&testCounter, nil},
			http.StatusOK,
			`{"delta":10,"id":"string","type":"counter"}`,
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
			body: "metric validator: type is invalid\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			storager := mocks.NewMockMetricStorager(ctrl)
			if tt.want != nil {
				storager.EXPECT().GetOne(gomock.Any(), tt.got.id, tt.got.Type).Return(tt.want.metric, tt.want.err)
			}

			apitest.New(tt.name).
				Handler(endpoints.NewAPIRouter(storager)).
				Post("/value/").
				Bodyf(`{"type":"%s","id":"%s"}`, tt.got.Type, tt.got.id).
				Expect(t).
				Status(tt.code).
				Body(tt.body).
				End()
		})
	}
}

func TestGetOneMetricFromJSON__contract_error(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	storager := mocks.NewMockMetricStorager(ctrl)
	storager.EXPECT().GetOne(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)

	apitest.Handler(endpoints.NewAPIRouter(storager)).
		Post("/value/").
		Body(`{"type}`).
		Expect(t).
		Status(http.StatusBadRequest).
		Body("parse error: unterminated string literal near offset 7 of '{\"type}'\n").
		End()
}
