package endpoints

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/novoseltcev/go-course/internal/model"
	"github.com/novoseltcev/go-course/internal/server/storage"
	"github.com/novoseltcev/go-course/internal/server/storage/mem"
)


func testRequest(t *testing.T, ts *httptest.Server, method, path string) (*http.Response, string) {
	req, err := http.NewRequest(method, ts.URL + path, http.NoBody)
	require.NoError(t, err)

	resp, err := ts.Client().Do(req)
	require.NoError(t, err)

	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	return resp, string(respBody)
}

func getRouter(s storage.MetricStorager) http.Handler {
	r := chi.NewRouter()
	r.Post(`/update/{metricType}/{metricName}/{metricValue}`, UpdateMetric(s))
	return r
}

func TestUpdateMetric(t *testing.T) {
	var counterValue = int64(1)
	var gaugeValue = float64(123.56)
	var newCounterValue = int64(3)
	var newGaugeValue = float64(234.)
	sharedStorage := mem.New()

	tests := []struct {
		name string
		storage storage.MetricStorager
		method string
		url string
		status int
		want *model.Metric
		length *int
	}{
		{
			name: "add new counter",
			storage: sharedStorage,
			method: http.MethodPost,
			url: "/update/counter/some/1",
			status: http.StatusOK,
			want: &model.Metric{Name: "some", Type: "counter", Delta: &counterValue},
		},
		{
			name: "add new gauge",
			storage: sharedStorage,
			method: http.MethodPost,
			url: "/update/gauge/some/123.56",
			status: http.StatusOK,
			want: &model.Metric{Name: "some", Type: "gauge", Value: &gaugeValue},
		},
		{
			name: "add exists counter",
			storage: sharedStorage,
			method: http.MethodPost,
			url: "/update/counter/some/2",
			status: http.StatusOK,
			want: &model.Metric{Name: "some", Type: "counter", Delta: &newCounterValue},
		},
		{
			name: "add exists gauge",
			storage: sharedStorage,
			method: http.MethodPost,
			url: "/update/gauge/some/234",
			status: http.StatusOK,
			want: &model.Metric{Name: "some", Type: "gauge", Value: &newGaugeValue},
		},
		{
			name: "invalid method",
			method: http.MethodGet,
			url: "/update/gauge/some/123.56",
			status: http.StatusMethodNotAllowed,
		},
		{
			name: "miss gauge value",
			method: http.MethodPost,
			url: "/update/gauge/some",
			status: http.StatusNotFound,
		},
		{
			name: "miss counter value",
			method: http.MethodPost,
			url: "/update/counter/some",
			status: http.StatusNotFound,
		},
		{
			name: "unknown metric type",
			method: http.MethodPost,
			url: "/update/some/some/1",
			status: http.StatusBadRequest,
		},
		{
			name: "invalid gauge value",
			method: http.MethodPost,
			url: "/update/gauge/some/value",
			status: http.StatusBadRequest,
		},
		{
			name: "invalid counter value",
			method: http.MethodPost,
			url: "/update/counter/some/1.",
			status: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			
			ts := httptest.NewServer(getRouter(tt.storage))
    		defer ts.Close()

			response, _ := testRequest(t, ts, tt.method, tt.url)
			defer response.Body.Close()

			assert.Equal(t, tt.status, response.StatusCode)

			if tt.want != nil {
				require.Equal(t, http.StatusOK, response.StatusCode)
				metric := *tt.want
				m, err := tt.storage.GetByName(context.TODO(), metric.Name, metric.Type)
				require.NoError(t, err, "Ошибка получения метрики из хранилища")
				require.NotNil(t, m, "Метрика не найдена")
				assert.Equal(t, metric, *m, "Отправленная и полученная метрики не равны")
			}
		})
	}
}
