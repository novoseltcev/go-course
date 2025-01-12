package collectors

import (
	"context"
	"math/rand"
	"runtime"

	"github.com/novoseltcev/go-course/internal/schemas"
)

func CollectRuntimeMetrics(_ context.Context) ([]schemas.Metric, error) {
	rtm := new(runtime.MemStats)
	runtime.ReadMemStats(rtm)

	metrics := []struct {
		string
		uint64
	}{
		{"Alloc", rtm.Alloc},
		{"Alloc", rtm.Alloc},
		{"BuckHashSys", rtm.BuckHashSys},
		{"Frees", rtm.Frees},
		{"GCSys", rtm.GCSys},
		{"HeapAlloc", rtm.HeapAlloc},
		{"HeapIdle", rtm.HeapIdle},
		{"HeapInuse", rtm.HeapInuse},
		{"HeapObjects", rtm.HeapObjects},
		{"HeapReleased", rtm.HeapReleased},
		{"HeapSys", rtm.HeapSys},
		{"LastGC", rtm.LastGC},
		{"Lookups", rtm.Lookups},
		{"MCacheInuse", rtm.MCacheInuse},
		{"MCacheSys", rtm.MCacheSys},
		{"MSpanInuse", rtm.MSpanInuse},
		{"MSpanSys", rtm.MSpanSys},
		{"Mallocs", rtm.Mallocs},
		{"NextGC", rtm.NextGC},
		{"NumForcedGC", uint64(rtm.NumForcedGC)},
		{"NumGC", uint64(rtm.NumGC)},
		{"OtherSys", rtm.OtherSys},
		{"PauseTotalNs", rtm.PauseTotalNs},
		{"StackInuse", rtm.StackInuse},
		{"StackSys", rtm.StackSys},
		{"Sys", rtm.Sys},
		{"TotalAlloc", rtm.TotalAlloc},
	}

	result := make([]schemas.Metric, 0, len(metrics)+2) // nolint:mnd

	for _, m := range metrics {
		val := float64(m.uint64)
		result = append(result, schemas.Metric{ID: m.string, MType: schemas.Gauge, Value: &val})
	}

	randValue := rand.Float64() // nolint:gosec
	result = append(result, schemas.Metric{ID: "RandomValue", MType: schemas.Gauge, Value: &randValue})

	var counterStep int64 = 10
	result = append(result, schemas.Metric{ID: "PollCount", MType: schemas.Counter, Delta: &counterStep})

	return result, nil
}
