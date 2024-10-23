package endpoints_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"

	"github.com/novoseltcev/go-course/internal/schemas"
	"github.com/novoseltcev/go-course/internal/server/endpoints"
	"github.com/novoseltcev/go-course/internal/storages"
)

func TestGetOneMetric(t *testing.T) {
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

	router := chi.NewRouter()
	router.Get(`/{metricType}/{metricName}`, endpoints.GetOneMetric(storage))
	t.Run("should return gauge metric value", func(t *testing.T) {
		t.Parallel()

		ts := httptest.NewServer(router)

		defer ts.Close()

		resp, body := testRequest(t, ts, http.MethodGet, "/gauge/test", http.NoBody)

		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Equal(t, strings.TrimRight(fmt.Sprintf("%f", testValue), "0"), body)
	})

	t.Run("should return counter metric value", func(t *testing.T) {
		t.Parallel()

		ts := httptest.NewServer(router)

		defer ts.Close()

		resp, body := testRequest(t, ts, http.MethodGet, "/counter/test", http.NoBody)

		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Equal(t, strconv.FormatInt(testDelta, 10), body)
	})

	t.Run("should return error if metric not found", func(t *testing.T) {
		t.Parallel()

		ts := httptest.NewServer(router)

		defer ts.Close()

		resp, body := testRequest(t, ts, http.MethodGet, "/gauge/unknown", http.NoBody)

		defer resp.Body.Close()

		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
		assert.Equal(t, "metric not found\n", body)
	})

	t.Run("should return error if metric type is invalid", func(t *testing.T) {
		t.Parallel()

		ts := httptest.NewServer(router)

		defer ts.Close()

		resp, body := testRequest(t, ts, http.MethodGet, "/unknown/test", http.NoBody)

		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		assert.Equal(t, "invalid metric type\n", body)
	})
}
