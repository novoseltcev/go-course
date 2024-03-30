package workers

import (
	"fmt"
	"math/rand"
	"runtime"

	"github.com/novoseltcev/go-course/internal/types"
)


func CollectMetrics(counterStorage *map[string]types.Counter, gaugeStorage *map[string]types.Gauge) func() {
	fmt.Println("init CollectMetrics worker")
	return func () {
		for k, v := range collectRuntimeMetrics() {
			(*gaugeStorage)[k] = types.Gauge(v)
		}
		(*counterStorage)["PollCount"] += 1
		(*gaugeStorage)["RandomValue"] = types.Gauge(rand.Float64())
	}
}

func collectRuntimeMetrics () map[string]float64 {
	rtm := new(runtime.MemStats)
	runtime.ReadMemStats(rtm)
	result := make(map[string]float64, 27)

	result["GCCPUFraction"] = rtm.GCCPUFraction
	result["Alloc"] = float64(rtm.Alloc)
	result["BuckHashSys"] = float64(rtm.BuckHashSys)
	result["Frees"] = float64(rtm.Frees)
	result["GCSys"] = float64(rtm.GCSys)
	result["HeapAlloc"] = float64(rtm.HeapAlloc)
	result["HeapIdle"] = float64(rtm.HeapIdle)
	result["HeapInuse"] = float64(rtm.HeapInuse)
	result["HeapObjects"] = float64(rtm.HeapObjects)
	result["HeapReleased"] = float64(rtm.HeapReleased)
	result["HeapSys"] = float64(rtm.HeapSys)
	result["LastGC"] = float64(rtm.LastGC)
	result["Lookups"] = float64(rtm.Lookups)
	result["MCacheInuse"] = float64(rtm.MCacheInuse)
	result["MCacheSys"] = float64(rtm.MCacheSys)
	result["MSpanInuse"] = float64(rtm.MSpanInuse)
	result["MSpanSys"] = float64(rtm.MSpanSys)
	result["Mallocs"] = float64(rtm.Mallocs)
	result["NextGC"] = float64(rtm.NextGC)
	result["NumForcedGC"] = float64(rtm.NumForcedGC)
	result["NumGC"] = float64(rtm.NumGC)
	result["OtherSys"] = float64(rtm.OtherSys)
	result["PauseTotalNs"] = float64(rtm.PauseTotalNs)
	result["StackInuse"] = float64(rtm.StackInuse)
	result["StackSys"] = float64(rtm.StackSys)
	result["Sys"] = float64(rtm.Sys)
	result["TotalAlloc"] = float64(rtm.TotalAlloc)

	return result
}
