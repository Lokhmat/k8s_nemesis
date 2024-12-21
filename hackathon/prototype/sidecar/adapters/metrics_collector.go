package adapters

import (
	"os"

	"sidecar/domain"
	"sidecar/port"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
)

type MetricsCollectorAdapter struct{}

func NewMetricsCollectorAdapter() port.MetricsCollector {
	return &MetricsCollectorAdapter{}
}

func (m *MetricsCollectorAdapter) Collect() (*domain.Metrics, error) {
	cpuPercent, err := cpu.Percent(0, false)
	if err != nil {
		return nil, err
	}

	memStats, err := mem.VirtualMemory()
	if err != nil {
		return nil, err
	}

	return &domain.Metrics{
		PodName:        os.Getenv("HOSTNAME"),
		DeploymentName: os.Getenv("DEPLOYMENT_NAME"),
		PodID:          os.Getenv("POD_ID"),
		CPUUsage:       cpuPercent[0],
		RAMUsage:       memStats.UsedPercent,
	}, nil
}
