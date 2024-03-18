package endpoints

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	. "github.com/novoseltcev/go-course/internal/server/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

type want[T Counter | Gauge] struct {
	metric Metric[T]
	length int
}

func TestUpdateMetric(t *testing.T) {
	type args struct {
		counterStorage MemStorage[Counter]
		gaugeStorage MemStorage[Gauge]
	}
	type result struct {
		counter *want[Counter]
		gauge *want[Gauge]
	}
	tests := []struct {
		name string
		args args
		method string
		url string
		status int
		result result
	}{
		{
			name: "add new counter",
			args: args{
				counterStorage: MemStorage[Counter]{
					Metrics: make(map[string]Counter),
				},
				gaugeStorage: MemStorage[Gauge]{
					Metrics: make(map[string]Gauge),
				},
			},
			method: http.MethodPost,
			url: "/update/counter/some/1",
			status: http.StatusOK,
			result: result{
				counter: &want[Counter]{
					metric: Metric[Counter]{Name: "some", Value: 1},
					length: 1,
				},
			},
		},
		{
			name: "add new gauge",
			args: args{
				counterStorage: MemStorage[Counter]{
					Metrics: make(map[string]Counter),
				},
				gaugeStorage: MemStorage[Gauge]{
					Metrics: make(map[string]Gauge),
				},
			},
			method: http.MethodPost,
			url: "/update/gauge/some/123.56",
			status: http.StatusOK,
			result: result{
				gauge: &want[Gauge]{
					metric: Metric[Gauge]{Name: "some", Value: 123.56},
					length: 1,
				},
			},
		},
		{
			name: "add exists counter",
			args: args{
				counterStorage: MemStorage[Counter]{
					Metrics: map[string]Counter{"some": 1},
				},
				gaugeStorage: MemStorage[Gauge]{
					Metrics: make(map[string]Gauge),
				},
			},
			method: http.MethodPost,
			url: "/update/counter/some/2",
			status: http.StatusOK,
			result: result{
				counter: &want[Counter]{
					metric: Metric[Counter]{Name: "some", Value: 3},
					length: 1,
				},
			},
		},
		{
			name: "add exists gauge",
			args: args{
				counterStorage: MemStorage[Counter]{
					Metrics: make(map[string]Counter),
				},
				gaugeStorage: MemStorage[Gauge]{
					Metrics: map[string]Gauge{"some": 11.32},
				},
			},
			method: http.MethodPost,
			url: "/update/gauge/some/123.56",
			status: http.StatusOK,
			result: result{
				gauge: &want[Gauge]{
					metric: Metric[Gauge]{Name: "some", Value: 123.56},
					length: 1,
				},
			},
		},
		{
			name: "invalid method",
			args: args{
				counterStorage: MemStorage[Counter]{Metrics: make(map[string]Counter)},
				gaugeStorage: MemStorage[Gauge]{Metrics: make(map[string]Gauge)},
			},
			method: http.MethodGet,
			url: "/update/gauge/some/123.56",
			status: http.StatusMethodNotAllowed,
			result: result{},
		},
		{
			name: "miss gauge value",
			args: args{
				counterStorage: MemStorage[Counter]{},
				gaugeStorage: MemStorage[Gauge]{},
			},
			method: http.MethodPost,
			url: "/update/gauge/some",
			status: http.StatusNotFound,
			result: result{},
		},
		{
			name: "miss counter value",
			args: args{
				counterStorage: MemStorage[Counter]{},
				gaugeStorage: MemStorage[Gauge]{},
			},
			method: http.MethodPost,
			url: "/update/counter/some",
			status: http.StatusNotFound,
			result: result{},
		},
		{
			name: "unknown metric type",
			args: args{
				counterStorage: MemStorage[Counter]{},
				gaugeStorage: MemStorage[Gauge]{},
			},
			method: http.MethodPost,
			url: "/update/some/some/1",
			status: http.StatusBadRequest,
			result: result{},
		},
		{
			name: "invalid gauge value",
			args: args{
				counterStorage: MemStorage[Counter]{},
				gaugeStorage: MemStorage[Gauge]{},
			},
			method: http.MethodPost,
			url: "/update/gauge/some/value",
			status: http.StatusBadRequest,
			result: result{},
		},
		{
			name: "invalid counter value",
			args: args{
				counterStorage: MemStorage[Counter]{},
				gaugeStorage: MemStorage[Gauge]{},
			},
			method: http.MethodPost,
			url: "/update/counter/some/1.",
			status: http.StatusBadRequest,
			result: result{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			var counterStorage Storage[Counter] = tt.args.counterStorage
			var gaugeStorage Storage[Gauge] = tt.args.gaugeStorage

			ts := httptest.NewServer(MetricRouter(&counterStorage, &gaugeStorage))
    		defer ts.Close()

			response, _ := testRequest(t, ts, tt.method, tt.url)
			defer response.Body.Close()

			assert.Equal(t, tt.status, response.StatusCode)

			if tt.result.counter != nil {
				require.Equal(t, http.StatusOK, response.StatusCode)

				metrics := tt.args.counterStorage.Metrics
				require.Len(t, metrics, tt.result.counter.length)
				metric := tt.result.counter.metric
				require.Contains(t, metrics, metric.Name)
				assert.Equal(t, metric.Value, metrics[metric.Name])
			}
			if tt.result.gauge != nil {
				require.Equal(t, http.StatusOK, response.StatusCode)

				metrics := tt.args.gaugeStorage.Metrics
				require.Len(t, metrics, tt.result.gauge.length)
				metric := tt.result.gauge.metric
				require.Contains(t, metrics, metric.Name)
				assert.Equal(t, metric.Value, metrics[metric.Name])
			}
		})
	}
}
