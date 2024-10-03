package workers

import (
	"context"
	"math/rand"
	"runtime"
	"time"

	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/cpu"
	log "github.com/sirupsen/logrus"

	"github.com/novoseltcev/go-course/internal/model"
)

var counterStep int64 = 1
type pair struct {
	name string
	value float64
}

func CollectMetrics(ctx context.Context, delay time.Duration) <-chan model.Metric {
	ch := make(chan model.Metric)
	go func() {
		defer close(ch)
		log.WithField("workerName", "CollectMetrics").Info("start worker")

		for {
			for _, p := range getRuntimeMetrics() {
				log.WithFields(log.Fields{"name": p.name, "value": p.value}).Info("collect runtime metric")
				ch <- model.Metric{Name: p.name, Type: "gauge", Value: &p.value}
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

func getRuntimeMetrics() []pair {
	rtm := new(runtime.MemStats)
	runtime.ReadMemStats(rtm)

	result := make([]pair, 27)
	result = append(result, pair{"GCCPUFraction", rtm.GCCPUFraction})
	result = append(result, pair{"Alloc", float64(rtm.Alloc)})
	result = append(result, pair{"BuckHashSys", float64(rtm.BuckHashSys)})
	result = append(result, pair{"Frees", float64(rtm.Frees)})
	result = append(result, pair{"GCSys", float64(rtm.GCSys)})
	result = append(result, pair{"HeapAlloc", float64(rtm.HeapAlloc)})
	result = append(result, pair{"HeapIdle", float64(rtm.HeapIdle)})
	result = append(result, pair{"HeapInuse", float64(rtm.HeapInuse)})
	result = append(result, pair{"HeapObjects", float64(rtm.HeapObjects)})
	result = append(result, pair{"HeapReleased", float64(rtm.HeapReleased)})
	result = append(result, pair{"HeapSys", float64(rtm.HeapSys)})
	result = append(result, pair{"LastGC", float64(rtm.LastGC)})
	result = append(result, pair{"Lookups", float64(rtm.Lookups)})
	result = append(result, pair{"MCacheInuse", float64(rtm.MCacheInuse)})
	result = append(result, pair{"MCacheSys", float64(rtm.MCacheSys)})
	result = append(result, pair{"MSpanInuse", float64(rtm.MSpanInuse)})
	result = append(result, pair{"MSpanSys", float64(rtm.MSpanSys)})
	result = append(result, pair{"Mallocs", float64(rtm.Mallocs)})
	result = append(result, pair{"NextGC", float64(rtm.NextGC)})
	result = append(result, pair{"NumForcedGC", float64(rtm.NumForcedGC)})
	result = append(result, pair{"NumGC", float64(rtm.NumGC)})
	result = append(result, pair{"OtherSys", float64(rtm.OtherSys)})
	result = append(result, pair{"PauseTotalNs", float64(rtm.PauseTotalNs)})
	result = append(result, pair{"StackInuse", float64(rtm.StackInuse)})
	result = append(result, pair{"StackSys", float64(rtm.StackSys)})
	result = append(result, pair{"Sys", float64(rtm.Sys)})
	result = append(result, pair{"TotalAlloc", float64(rtm.TotalAlloc)})

	return result
}

func getCoreMetrics() []pair {
	vmStat, _ := mem.VirtualMemory()
	cpuLoad, _ := cpu.Percent(0, true)

	result := make([]pair, 3)
	result = append(result, pair{"TotalMemory", float64(vmStat.Total)})
	result = append(result, pair{"FreeMemory", float64(vmStat.Free)})
	result = append(result, pair{"CPUutilization1", cpuLoad[0]})

	return result
}

func CollectCoreMetrics(ctx context.Context, delay time.Duration) <-chan model.Metric {
	ch := make(chan model.Metric)
	go func() {
		defer close(ch)
		log.WithField("workerName", "CollectAdditionalMetrics").Info("start worker")

		for {
			for _, p := range getCoreMetrics() {
				log.WithFields(log.Fields{"name": p.name, "value": p.value}).Info("collect core metric")
				ch <- model.Metric{Name: p.name, Type: "gauge", Value: &p.value}
			}

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
