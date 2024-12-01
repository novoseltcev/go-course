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

func TestUpdateBatch(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	type got struct {
		body  string
		batch []schemas.Metric
	}

	tests := []struct {
		name string
		got
		err  error
		code int
		body string
	}{
		{
			name: "success",
			got: got{
				body:  `[{"id":"test","type":"gauge","value":10.123},{"id":"test","type":"counter","delta":10}]`,
				batch: []schemas.Metric{testGauge, testCounter},
			},
			code: http.StatusOK,
		},
		{
			name: "skip one",
			got: got{
				body:  `[{"id":"test","type":"gauge","value":10.123},{"id":"test","type":"unknown","value":10.123}]`,
				batch: []schemas.Metric{testGauge},
			},
			code: http.StatusOK,
		},
		{
			name: "skip all",
			got: got{
				body:  `[{"id":"test","type":"unknown","value":10.123}]`,
				batch: nil,
			},
			code: http.StatusOK,
		},
		{
			name: "failed save",
			got: got{
				body:  `[{"id":"test","type":"gauge","value":10.123},{"id":"test","type":"gauge","value":10.123}]`,
				batch: []schemas.Metric{testGauge, testGauge},
			},
			err:  errTest,
			code: http.StatusInternalServerError,
			body: "failed to save metrics\n",
		},
		{
			name: "unmarshalable body",
			got:  got{body: `[{"unknown"}]`},
			code: http.StatusBadRequest,
			body: "parse error: syntax error near offset 11 of '[{\"unknown\"}]'\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			storager := mocks.NewMockMetricStorager(ctrl)
			router := endpoints.NewAPIRouter(storager)

			if tt.got.batch != nil {
				storager.EXPECT().SaveBatch(gomock.Any(), tt.got.batch).Times(1).Return(tt.err)
			} else {
				storager.EXPECT().SaveBatch(gomock.Any(), gomock.Any()).Times(0)
			}

			req := httptest.NewRequest(http.MethodPost, "/updates/", bytes.NewBufferString(tt.got.body))
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.code, w.Code)
			assert.Equal(t, tt.body, w.Body.String())
		})
	}
}
