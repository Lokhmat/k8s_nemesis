package repositories

import (
	"autoscaler/internal/domain/entities"
	"context"
)

type MappingRepository interface {
	GetDeploymentByAPIPrefix(ctx context.Context, apiPrefix string) (*entities.Deployment, error)
	LoadMappings(ctx context.Context) error
	SubscribeToUpdates(ctx context.Context)
}
