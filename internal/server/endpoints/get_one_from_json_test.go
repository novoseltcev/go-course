package endpoints_test

import (
	"bytes"
	"fmt"
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
			got{id: testID, Type: schemas.Gauge},
			&want{&testGauge, nil},
			http.StatusOK,
			`{"value":10.123,"id":"test","type":"gauge"}`,
		},
		{
			"success counter",
			got{id: testID, Type: schemas.Counter},
			&want{&testCounter, nil},
			http.StatusOK,
			`{"delta":10,"id":"test","type":"counter"}`,
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			storager := mocks.NewMockMetricStorager(ctrl)
			router := endpoints.NewAPIRouter(storager)

			if tt.want != nil {
				storager.EXPECT().GetOne(gomock.Any(), tt.got.id, tt.got.Type).Return(tt.want.metric, tt.want.err)
			}

			req := httptest.NewRequest(http.MethodPost, "/value/", bytes.NewBufferString(
				fmt.Sprintf(`{"type":"%s","id":"%s"}`, tt.got.Type, tt.got.id),
			))
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.code, w.Code)
			assert.Equal(t, tt.body, w.Body.String())
		})
	}
}

func TestGetOneMetricFromJSON__contract_error(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	storager := mocks.NewMockMetricStorager(ctrl)
	router := endpoints.NewAPIRouter(storager)

	storager.EXPECT().GetOne(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)

	req := httptest.NewRequest(http.MethodPost, "/value/", bytes.NewBufferString(`{"type}`))
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, "parse error: unterminated string literal near offset 7 of '{\"type}'\n", w.Body.String())
}
