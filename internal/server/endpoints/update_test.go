package endpoints_test

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/novoseltcev/go-course/internal/schemas"
	"github.com/novoseltcev/go-course/internal/server/endpoints"
	"github.com/novoseltcev/go-course/internal/storages"
)

func testRequest(t *testing.T, ts *httptest.Server, method, path string, body io.Reader) (*http.Response, string) {
	t.Helper()

	req, err := http.NewRequestWithContext(context.Background(), method, ts.URL+path, body)
	require.NoError(t, err)

	resp, err := ts.Client().Do(req)
	require.NoError(t, err)

	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	return resp, string(respBody)
}

func TestUpdateMetric(t *testing.T) { //nolint:paralleltest
	var (
		testCounterValue    = int64(1)
		testGaugeValue      = float64(123.56)
		testNewCounterValue = int64(3)
		testNewGaugeValue   = float64(234.)
		sharedStorage       = storages.NewMemStorage()
	)

	tests := []struct {
		name    string
		storage storages.MetricStorager
		method  string
		url     string
		status  int
		want    *schemas.Metric
	}{
		{
			name:    "add new counter",
			storage: sharedStorage,
			method:  http.MethodPost,
			url:     "/counter/some/1",
			status:  http.StatusOK,
			want:    &schemas.Metric{ID: "some", MType: schemas.Counter, Delta: &testCounterValue},
		},
		{
			name:    "add new gauge",
			storage: sharedStorage,
			method:  http.MethodPost,
			url:     "/gauge/some/123.56",
			status:  http.StatusOK,
			want:    &schemas.Metric{ID: "some", MType: schemas.Gauge, Value: &testGaugeValue},
		},
		{
			name:    "add exists counter",
			storage: sharedStorage,
			method:  http.MethodPost,
			url:     "/counter/some/2",
			status:  http.StatusOK,
			want:    &schemas.Metric{ID: "some", MType: schemas.Counter, Delta: &testNewCounterValue},
		},
		{
			name:    "add exists gauge",
			storage: sharedStorage,
			method:  http.MethodPost,
			url:     "/gauge/some/234",
			status:  http.StatusOK,
			want:    &schemas.Metric{ID: "some", MType: schemas.Gauge, Value: &testNewGaugeValue},
		},
		{
			name:   "invalid method",
			method: http.MethodGet,
			url:    "/gauge/some/123.56",
			status: http.StatusMethodNotAllowed,
		},
		{
			name:   "miss gauge value",
			method: http.MethodPost,
			url:    "/gauge/some",
			status: http.StatusNotFound,
		},
		{
			name:   "miss counter value",
			method: http.MethodPost,
			url:    "/counter/some",
			status: http.StatusNotFound,
		},
		{
			name:   "unknown metric type",
			method: http.MethodPost,
			url:    "/some/some/1",
			status: http.StatusBadRequest,
		},
		{
			name:   "invalid gauge value",
			method: http.MethodPost,
			url:    "/gauge/some/value",
			status: http.StatusBadRequest,
		},
		{
			name:   "invalid counter value",
			method: http.MethodPost,
			url:    "/counter/some/1.",
			status: http.StatusBadRequest,
		},
	}
	for _, tt := range tests { //nolint:paralleltest
		t.Run(tt.name, func(t *testing.T) {
			router := chi.NewRouter()
			router.Post(`/{metricType}/{metricName}/{metricValue}`, endpoints.UpdateMetric(tt.storage))

			ts := httptest.NewServer(router)

			defer ts.Close()

			resp, _ := testRequest(t, ts, tt.method, tt.url, http.NoBody)

			defer resp.Body.Close()

			assert.Equal(t, tt.status, resp.StatusCode)

			if tt.want != nil {
				require.Equal(t, http.StatusOK, resp.StatusCode)

				metric := *tt.want
				m, err := tt.storage.GetByName(context.TODO(), metric.ID, metric.MType)
				require.NoError(t, err, "Ошибка получения метрики из хранилища")
				require.NotNil(t, m, "Метрика не найдена")
				assert.Equal(t, metric, *m, "Отправленная и полученная метрики не равны")
			}
		})
	}
}
