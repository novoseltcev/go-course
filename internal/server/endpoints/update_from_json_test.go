package endpoints_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	gomock "go.uber.org/mock/gomock"

	"github.com/novoseltcev/go-course/internal/schemas"
	"github.com/novoseltcev/go-course/internal/server/endpoints"
	"github.com/novoseltcev/go-course/mocks"
)

func TestUpdateMetricFromJSON(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	type got struct {
		body   string
		metric *schemas.Metric
	}

	tests := []struct {
		name string
		got
		err  error
		code int
		body string
	}{
		{
			name: "success gauge",
			got:  got{body: `{"id":"test","type":"gauge","value":10.123}`, metric: &testGauge},
			code: http.StatusOK,
		},
		{
			name: "success counter",
			got:  got{body: `{"id":"test","type":"counter","delta":10}`, metric: &testCounter},
			code: http.StatusOK,
		},
		{
			name: "invalid metric type",
			got:  got{body: `{"id":"test","type":"unknown","value":10.123}`},
			code: http.StatusBadRequest,
			body: "type is invalid\n",
		},
		{
			name: "invalid gauge value",
			got:  got{body: `{"id":"test","type":"gauge","value":"value"}`},
			code: http.StatusBadRequest,
			body: "parse error: expected number near offset 43 of 'value'\n",
		},
		{
			name: "invalid counter value",
			got:  got{body: `{"id":"test","type":"counter","delta":"1."}`},
			code: http.StatusBadRequest,
			body: "parse error: expected number near offset 42 of '1.'\n",
		},
		{
			name: "failed save",
			got:  got{body: `{"id":"test","type":"gauge","value":10.123}`, metric: &testGauge},
			err:  errTest,
			code: http.StatusInternalServerError,
			body: "failed to save metric\n",
		},
		{
			name: "invalid contract",
			got:  got{body: `{}`},
			code: http.StatusBadRequest,
			body: "id is required\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			storager := mocks.NewMockMetricStorager(ctrl)
			router := endpoints.NewAPIRouter(storager)

			if tt.got.metric != nil {
				storager.EXPECT().Save(gomock.Any(), tt.got.metric).Return(tt.err)
			}

			req := httptest.NewRequest(http.MethodPost, "/update/", bytes.NewBufferString(tt.got.body))
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.code, w.Code)
			assert.Equal(t, tt.body, w.Body.String())
		})
	}
}
