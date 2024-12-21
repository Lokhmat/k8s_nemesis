package application

import (
	"fmt"
	"sidecar/port"
	"time"
)

type Sidecar struct {
	collector port.MetricsCollector
	pusher    port.MetricsPusher
	interval  time.Duration
}

func NewSidecar(collector port.MetricsCollector, pusher port.MetricsPusher, interval time.Duration) *Sidecar {
	return &Sidecar{
		collector: collector,
		pusher:    pusher,
		interval:  interval,
	}
}

func (s *Sidecar) Run() {
	for {
		metrics, err := s.collector.Collect()
		if err != nil {
			fmt.Printf("Error collecting metrics: %v\n", err)
			time.Sleep(s.interval)
			continue
		}

		fmt.Printf("Collected Metrics: %+v\n", *metrics)

		if err := s.pusher.Push(metrics); err != nil {
			fmt.Printf("Error pushing metrics: %v\n", err)
		}

		time.Sleep(s.interval)
	}
}
