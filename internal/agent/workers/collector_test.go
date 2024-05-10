package workers

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)


func TestCollectMetricsToEmptyStorage(t *testing.T) {
	counterStorage := make(map[string]int64)
	gaugeStorage := make(map[string]float64)

	require.Empty(t, gaugeStorage)
	require.Empty(t, counterStorage)

	CollectMetrics(&counterStorage, &gaugeStorage)

	require.Len(t, gaugeStorage, 28)
	gaugeStorageKeys := make([]string, 0, 28)
	for key := range gaugeStorage {
		gaugeStorageKeys = append(gaugeStorageKeys, key)
	}

	assert.ElementsMatch(t, []string{
		"Alloc",
		"BuckHashSys",
		"Frees",
		"GCCPUFraction",
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
	}, gaugeStorageKeys)


	require.Len(t, counterStorage, 1)
	require.Contains(t, counterStorage, "PollCount")
	assert.Equal(t, int64(1), counterStorage["PollCount"])
}


func TestCollectMetricsToFullStorage(t *testing.T) {
	var pollCount int64 = 1
	counterStorage := map[string]int64{"PollCount": pollCount}
	gaugeStorage := map[string]float64{
		"Alloc": 0.0,
		"BuckHashSys": 0.0,
		"Frees": 0.0,
		"GCCPUFraction": 0.0,
		"GCSys": 0.0,
		"HeapAlloc": 0.0,
		"HeapIdle": 0.0,
		"HeapInuse": 0.0,
		"HeapObjects": 0.0,
		"HeapReleased": 0.0,
		"HeapSys": 0.0,
		"LastGC": 0.0,
		"Lookups": 0.0,
		"MCacheInuse": 0.0,
		"MCacheSys": 0.0,
		"MSpanInuse": 0.0,
		"MSpanSys": 0.0,
		"Mallocs": 0.0,
		"NextGC": 0.0,
		"NumForcedGC": 0.0,
		"NumGC": 0.0,
		"OtherSys": 0.0,
		"PauseTotalNs": 0.0,
		"StackInuse": 0.0,
		"StackSys": 0.0,
		"Sys": 0.0,
		"TotalAlloc": 0.0,
		"RandomValue": 0.0,
	}

	assert.Len(t, gaugeStorage, 28)
	assert.Len(t, counterStorage, 1)

	CollectMetrics(&counterStorage, &gaugeStorage)

	assert.Len(t, gaugeStorage, 28)

	gaugeStorageKeys := make([]string, 0, 28)
	for key := range gaugeStorage {
		gaugeStorageKeys = append(gaugeStorageKeys, key)
	}

	assert.ElementsMatch(t, []string{
		"Alloc",
		"BuckHashSys",
		"Frees",
		"GCCPUFraction",
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
	}, gaugeStorageKeys)


	assert.Len(t, counterStorage, 1)
	assert.Contains(t, counterStorage, "PollCount")
	assert.Equal(t, pollCount + 1, counterStorage["PollCount"])
}
