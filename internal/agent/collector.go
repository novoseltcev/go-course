package agent

import (
	"math/rand"
	"fmt"
	"runtime"
)

type Storage[T int64 | float64] map[string]T

func CollectMetrics(counterStorage *Storage[int64], gaugeStorage *Storage[float64]) func() {
	fmt.Println("init CollectMetrics worker")
	return func () {
		collectRuntimeMetrics(gaugeStorage)
		(*counterStorage)["PollCount"] += 1
		(*gaugeStorage)["RandomValue"] = rand.Float64()
	}
}

func collectRuntimeMetrics (storage *Storage[float64]) {
	rtm := new(runtime.MemStats)
	runtime.ReadMemStats(rtm)
	(*storage)["Alloc"] = float64(rtm.Alloc)
	(*storage)["BuckHashSys"] = float64(rtm.BuckHashSys)
	(*storage)["Frees"] = float64(rtm.Frees)
	(*storage)["GCCPUFraction"] = float64(rtm.GCCPUFraction)
	(*storage)["GCSys"] = float64(rtm.GCSys)
	(*storage)["HeapAlloc"] = float64(rtm.HeapAlloc)
	(*storage)["HeapIdle"] = float64(rtm.HeapIdle)
	(*storage)["HeapInuse"] = float64(rtm.HeapInuse)
	(*storage)["HeapObjects"] = float64(rtm.HeapObjects)
	(*storage)["HeapReleased"] = float64(rtm.HeapReleased)
	(*storage)["HeapSys"] = float64(rtm.HeapSys)
	(*storage)["LastGC"] = float64(rtm.LastGC)
	(*storage)["Lookups"] = float64(rtm.Lookups)
	(*storage)["MCacheInuse"] = float64(rtm.MCacheInuse)
	(*storage)["MCacheSys"] = float64(rtm.MCacheSys)
	(*storage)["MSpanInuse"] = float64(rtm.MSpanInuse)
	(*storage)["MSpanSys"] = float64(rtm.MSpanSys)
	(*storage)["Mallocs"] = float64(rtm.Mallocs)
	(*storage)["NextGC"] = float64(rtm.NextGC)
	(*storage)["NumForcedGC"] = float64(rtm.NumForcedGC)
	(*storage)["NumGC"] = float64(rtm.NumGC)
	(*storage)["OtherSys"] = float64(rtm.OtherSys)
	(*storage)["PauseTotalNs"] = float64(rtm.PauseTotalNs)
	(*storage)["StackInuse"] = float64(rtm.StackInuse)
	(*storage)["StackSys"] = float64(rtm.StackSys)
	(*storage)["Sys"] = float64(rtm.Sys)
	(*storage)["TotalAlloc"] = float64(rtm.TotalAlloc)
}
