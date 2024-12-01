package application

import (
	httpAdapter "autoscaler/internal/adapters/http"
	"autoscaler/internal/usecase"
	"context"
)

func (app *Application) initUseCases(_ context.Context) error {
	proxyFunc := httpAdapter.ProxyRequest(app.metricsRepo)
	handleRequestUseCase := usecase.NewHandleRequestUseCase(
		app.mappingRepo,
		app.deploymentClient,
		app.metricsRepo,
		proxyFunc,
	)
	app.server = httpAdapter.NewServer(handleRequestUseCase)
	return nil
}
