package collectors

import (
	"context"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"

	"github.com/novoseltcev/go-course/internal/schemas"
)

func CollectCoreMetrics(ctx context.Context) ([]schemas.Metric, error) {
	vmStat, err := mem.VirtualMemoryWithContext(ctx)
	if err != nil {
		return nil, err
	}

	cpuLoad, err := cpu.PercentWithContext(ctx, 0, true)
	if err != nil {
		return nil, err
	}

	total := float64(vmStat.Total)
	free := float64(vmStat.Free)

	return []schemas.Metric{
		{ID: "TotalMemory", MType: schemas.Gauge, Value: &total},
		{ID: "FreeMemory", MType: schemas.Gauge, Value: &free},
		{ID: "CPUutilization1", MType: schemas.Gauge, Value: &cpuLoad[0]},
	}, nil
}
