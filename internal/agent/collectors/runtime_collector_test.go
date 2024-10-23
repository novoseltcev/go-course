package collectors_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/novoseltcev/go-course/internal/agent/collectors"
	"github.com/novoseltcev/go-course/internal/schemas"
)

func TestCollectRuntimeMetrics(t *testing.T) {
	t.Parallel()

	res, err := collectors.CollectRuntimeMetrics(context.TODO())

	require.NoError(t, err)
	assert.Len(t, res, 29)

	keys := make([]string, 0, 29)

	for i, metric := range res {
		keys = append(keys, metric.ID)

		if i == 28 {
			assert.Equal(t, "PollCount", metric.ID)
			assert.Equal(t, schemas.Counter, metric.MType)
			assert.NotNil(t, metric.Delta)
			assert.Nil(t, metric.Value)

			continue
		}

		assert.Equal(t, schemas.Gauge, metric.MType)
		assert.NotNil(t, metric.Value)
		assert.Nil(t, metric.Delta)
	}

	assert.EqualValues(t, []string{
		"Alloc",
		"Alloc",
		"BuckHashSys",
		"Frees",
		"GCSys",
		"HeapAlloc",
		"HeapIdle",
		"HeapInuse",
		"HeapObjects",
		"HeapReleased",
		"HeapSys",
		"LastGC",
		"Lookups",
		"MCacheInuse",
		"MCacheSys",
		"MSpanInuse",
		"MSpanSys",
		"Mallocs",
		"NextGC",
		"NumForcedGC",
		"NumGC",
		"OtherSys",
		"PauseTotalNs",
		"StackInuse",
		"StackSys",
		"Sys",
		"TotalAlloc",
		"RandomValue",
		"PollCount",
	}, keys)
}
