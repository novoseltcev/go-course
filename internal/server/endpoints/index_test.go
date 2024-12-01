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

func TestIndex(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	type want struct {
		metrics []schemas.Metric
		err     error
	}

	tests := []struct {
		name string
		want
		code int
		body string
	}{
		{
			"success",
			want{[]schemas.Metric{testCounter, testGauge}, nil},
			http.StatusOK,
			`<!DOCTYPE html>
<html>
  <head>
    <meta charset="UTF-8">
    <title>Metrics</title>
  </head>
  <body>
    <div>
      <h2>Metrics</h2>
      <ul>
        <li>
          <div>
            <b>Counter test</b>: 10
          </div>
        </li>
        <li>
          <div>
            <b>Gauge test</b>: 10.123
          </div>
        </li>
      </ul>
    </div>
  </body>
</html>
`,
		},
		{
			"empty",
			want{[]schemas.Metric{}, nil},
			http.StatusOK,
			`<!DOCTYPE html>
<html>
  <head>
    <meta charset="UTF-8">
    <title>Metrics</title>
  </head>
  <body>
    <div>
      <h2>Metrics</h2>
      <div>Empty</div>
    </div>
  </body>
</html>
`,
		},
		{
			"failed get",
			want{nil, errTest},
			http.StatusInternalServerError,
			"failed to get metrics\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			storager := mocks.NewMockMetricStorager(ctrl)
			router := endpoints.NewAPIRouter(storager)

			storager.EXPECT().GetAll(gomock.Any()).Return(tt.want.metrics, tt.want.err)

			req := httptest.NewRequest(http.MethodGet, "/", http.NoBody)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.code, w.Code)
			assert.Equal(t, tt.body, w.Body.String())
		})
	}
}
