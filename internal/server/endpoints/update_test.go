package endpoints_test

import (
	"net/http"
	"testing"

	"github.com/steinfletcher/apitest"
	gomock "go.uber.org/mock/gomock"

	"github.com/novoseltcev/go-course/internal/schemas"
	"github.com/novoseltcev/go-course/internal/server/endpoints"
	"github.com/novoseltcev/go-course/mocks"
	"github.com/novoseltcev/go-course/pkg/testutils"
)

func TestUpdateMetric(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	tests := []struct {
		name   string
		url    string
		metric *schemas.Metric
		err    error
		code   int
		body   string
	}{
		{
			name:   "success counter",
			url:    "/update/counter/string/10",
			metric: &testCounter,
			code:   http.StatusOK,
		},
		{
			name:   "success gauge",
			url:    "/update/gauge/string/10.123",
			metric: &testGauge,
			code:   http.StatusOK,
		},
		{
			name:   "failed save",
			url:    "/update/gauge/string/10.123",
			metric: &testGauge,
			err:    testutils.Err,
			code:   http.StatusInternalServerError,
			body:   "failed to save metric\n",
		},
		{
			name: "unknown metric type",
			url:  "/update/unknown/string/10",
			code: http.StatusBadRequest,
			body: "metric validator: type is invalid\n",
		},
		{
			name: "invalid gauge value",
			url:  "/update/gauge/string/value",
			code: http.StatusBadRequest,
			body: "value is invalid\n",
		},
		{
			name: "invalid counter value",
			url:  "/update/counter/string/1.",
			code: http.StatusBadRequest,
			body: "delta is invalid\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			storager := mocks.NewMockMetricStorager(ctrl)
			if tt.metric != nil {
				storager.EXPECT().Save(gomock.Any(), tt.metric).Return(tt.err)
			}

			apitest.New(tt.name).
				Handler(endpoints.NewAPIRouter(storager)).
				Post(tt.url).
				Expect(t).
				Status(tt.code).
				Body(tt.body).
				End()
		})
	}
}
