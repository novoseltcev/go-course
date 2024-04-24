package endpoints

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

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


func TestUpdateMetric(t *testing.T) {
	metrics := make(map[string]map[string]model.Metric)
	metrics["counter"] = make(map[string]model.Metric)
	metrics["gauge"] = make(map[string]model.Metric)
	var counterValue int64 = 1
	var gaugeValue float64 = 123.56
	var newCounterValue int64 = 20
	var newGaugeValue float64 = 234.

	tests := []struct {
		name string
		storage mem.Storage
		method string
		url string
		status int
		want *model.Metric
		length *int
	}{
		{
			name: "add new counter",
			storage: mem.Storage{Metrics: metrics},
			method: http.MethodPost,
			url: "/update/counter/some/1",
			status: http.StatusOK,
			want: &model.Metric{Name: "some", Type: "counter", Delta: &counterValue},
		},
		{
			name: "add new gauge",
			storage: mem.Storage{Metrics: metrics},
			method: http.MethodPost,
			url: "/update/gauge/some/123.56",
			status: http.StatusOK,
			want: &model.Metric{Name: "some", Type: "gauge", Value: &gaugeValue},
		},
		{
			name: "add exists counter",
			storage: mem.Storage{Metrics: metrics},
			method: http.MethodPost,
			url: "/update/counter/some/20",
			status: http.StatusOK,
			want: &model.Metric{Name: "some", Type: "counter", Delta: &newCounterValue},
		},
		{
			name: "add exists gauge",
			storage: mem.Storage{Metrics: metrics},
			method: http.MethodPost,
			url: "/update/gauge/some/234",
			status: http.StatusOK,
			want: &model.Metric{Name: "some", Type: "gauge", Value: &newGaugeValue},
		},
		{
			name: "invalid method",
			storage: mem.Storage{},
			method: http.MethodGet,
			url: "/update/gauge/some/123.56",
			status: http.StatusMethodNotAllowed,
		},
		{
			name: "miss gauge value",
			storage: mem.Storage{},
			method: http.MethodPost,
			url: "/update/gauge/some",
			status: http.StatusNotFound,
		},
		{
			name: "miss counter value",
			storage: mem.Storage{},
			method: http.MethodPost,
			url: "/update/counter/some",
			status: http.StatusNotFound,
		},
		{
			name: "unknown metric type",
			storage: mem.Storage{},
			method: http.MethodPost,
			url: "/update/some/some/1",
			status: http.StatusBadRequest,
		},
		{
			name: "invalid gauge value",
			storage: mem.Storage{},
			method: http.MethodPost,
			url: "/update/gauge/some/value",
			status: http.StatusBadRequest,
		},
		{
			name: "invalid counter value",
			storage: mem.Storage{},
			method: http.MethodPost,
			url: "/update/counter/some/1.",
			status: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			var storage storage.MetricStorager = &tt.storage

			ts := httptest.NewServer(GetRouter(nil, &storage))
    		defer ts.Close()

			response, _ := testRequest(t, ts, tt.method, tt.url)
			defer response.Body.Close()

			assert.Equal(t, tt.status, response.StatusCode)

			if tt.want != nil {
				require.Equal(t, http.StatusOK, response.StatusCode)
				metric := *tt.want
				metrics := tt.storage.Metrics[metric.Type]
				require.Len(t, metrics, 1)
				require.Contains(t, metrics, metric.Name)
				assert.Equal(t, metric, metrics[metric.Name])
			}
		})
	}
}
