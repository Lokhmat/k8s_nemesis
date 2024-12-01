package application

import (
	"autoscaler/internal/adapters/kubernetes"
	"context"
)

func (app *Application) initServices(_ context.Context) error {
	app.deploymentClient = kubernetes.NewDeploymentClient(app.clientSet, app.config.Namespace)
	return nil
}
