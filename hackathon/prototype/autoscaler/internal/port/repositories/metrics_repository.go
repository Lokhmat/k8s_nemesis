package repositories

import "context"

type MetricsRepository interface {
	UpdateLastRequestTime(ctx context.Context, deploymentName string) error
	UpdateResponseTime(ctx context.Context, deploymentName string, responseTime int64) error
	Shutdown()
}
