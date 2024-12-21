package port

import "sidecar/domain"

type MetricsPusher interface {
	Push(metrics *domain.Metrics) error
}
