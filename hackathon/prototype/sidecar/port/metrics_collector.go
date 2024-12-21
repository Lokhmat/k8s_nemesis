package port

import "sidecar/domain"

type MetricsCollector interface {
	Collect() (*domain.Metrics, error)
}
