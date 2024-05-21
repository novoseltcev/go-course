package workers

import (
	"context"
	"math/rand"
	"runtime"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/novoseltcev/go-course/internal/model"
)

var counterStep int64 = 1

func CollectMetrics(ctx context.Context, delay time.Duration) <-chan model.Metric {
	ch := make(chan model.Metric)
	go func() {
		defer close(ch)
		log.WithField("workerName", "CollectMetrics").Info("start worker")

		for {
			for name, value := range collectRuntimeMetrics() {
				log.WithFields(log.Fields{"name": name, "value": value}).Info("collect runtime metric")
				ch <- model.Metric{Name: name, Type: "gauge", Value: &value}
			}

			randValue := rand.Float64()
			ch <- model.Metric{Name: "RandomValue", Type: "gauge", Value: &randValue}
			ch <- model.Metric{Name: "PollCount", Type: "counter", Delta: &counterStep}

			select {
			case <- ctx.Done():
				return
			default:
				time.Sleep(delay)
			}
		}
	}()
	return ch
}

func collectRuntimeMetrics() map[string]float64 {
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
