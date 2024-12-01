package services

import (
	"autoscaler/internal/domain/entities"
	"context"
)

type DeploymentClient interface {
	ScaleDeployment(ctx context.Context, deployment *entities.Deployment, replicas int32) error
	ArePodsAvailable(ctx context.Context, deployment *entities.Deployment) (bool, error)
}
