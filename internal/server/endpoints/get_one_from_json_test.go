package endpoints_test

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/novoseltcev/go-course/internal/schemas"
	"github.com/novoseltcev/go-course/internal/server/endpoints"
	"github.com/novoseltcev/go-course/internal/storages"
)

func TestGetOneMetricFromJSON(t *testing.T) {
	t.Parallel()

	var (
		testValue = 10.123
		testDelta = int64(10)
	)

	storage := storages.NewMemStorage()
	storage.SaveAll(context.TODO(), []schemas.Metric{
		{
			ID:    "test",
			MType: schemas.Gauge,
			Value: &testValue,
		},
		{
			ID:    "test",
			MType: schemas.Counter,
			Delta: &testDelta,
		},
	})
	t.Run("should return gauge metric value", func(t *testing.T) {
		t.Parallel()

		ts := httptest.NewServer(endpoints.GetOneMetricFromJSON(storage))

		defer ts.Close()

		resp, body := testRequest(
			t,
			ts,
			http.MethodPost,
			"/",
			bytes.NewBufferString(`{"type":"gauge","id":"test"}`),
		)

		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Equal(t, `{"id":"test","type":"gauge","value":10.123}`, body)
	})

	t.Run("should return counter metric value", func(t *testing.T) {
		t.Parallel()

		ts := httptest.NewServer(endpoints.GetOneMetricFromJSON(storage))

		defer ts.Close()

		resp, body := testRequest(
			t,
			ts,
			http.MethodPost,
			"/",
			bytes.NewBufferString(`{"type":"counter","id":"test"}`),
		)

		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Equal(t, `{"id":"test","type":"counter","delta":10}`, body)
	})

	t.Run("should return error if metric not found", func(t *testing.T) {
		t.Parallel()

		ts := httptest.NewServer(endpoints.GetOneMetricFromJSON(storage))

		defer ts.Close()

		resp, body := testRequest(
			t,
			ts,
			http.MethodPost,
			"/",
			bytes.NewBufferString(`{"type":"gauge","id":"unknown"}`),
		)

		defer resp.Body.Close()

		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
		assert.Equal(t, "metric not found\n", body)
	})

	t.Run("should return error if metric type is invalid", func(t *testing.T) {
		t.Parallel()

		ts := httptest.NewServer(endpoints.GetOneMetricFromJSON(storage))

		defer ts.Close()

		resp, body := testRequest(
			t,
			ts,
			http.MethodPost,
			"/",
			bytes.NewBufferString(`{"type":"unknown","id":"test"}`),
		)

		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		assert.Equal(t, "invalid metric type\n", body)
	})
}
