package endpoints_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	gomock "go.uber.org/mock/gomock"

	"github.com/novoseltcev/go-course/internal/schemas"
	"github.com/novoseltcev/go-course/internal/server/endpoints"
	"github.com/novoseltcev/go-course/mocks"
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
			url:    "/update/counter/test/10",
			metric: &testCounter,
			code:   http.StatusOK,
		},
		{
			name:   "success gauge",
			url:    "/update/gauge/test/10.123",
			metric: &testGauge,
			code:   http.StatusOK,
		},
		{
			name:   "failed save",
			url:    "/update/gauge/test/10.123",
			metric: &testGauge,
			err:    errTest,
			code:   http.StatusInternalServerError,
			body:   "failed to save metric\n",
		},
		{
			name: "unknown metric type",
			url:  "/update/unknown/test/10",
			code: http.StatusBadRequest,
			body: "type is invalid\n",
		},
		{
			name: "invalid gauge value",
			url:  "/update/gauge/test/value",
			code: http.StatusBadRequest,
			body: "invalid value\n",
		},
		{
			name: "invalid counter value",
			url:  "/update/counter/test/1.",
			code: http.StatusBadRequest,
			body: "invalid delta\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			storager := mocks.NewMockMetricStorager(ctrl)
			router := endpoints.NewAPIRouter(storager)

			if tt.metric != nil {
				storager.EXPECT().Save(gomock.Any(), tt.metric).Return(tt.err)
			}

			req := httptest.NewRequest(http.MethodPost, tt.url, http.NoBody)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.code, w.Code)
			assert.Equal(t, tt.body, w.Body.String())
		})
	}
}
