package application

import (
	redisAdapter "autoscaler/internal/adapters/redis"
	"context"
	"time"
)

func (app *Application) initRepositories(_ context.Context) error {
	var err error
	app.mappingRepo, err = redisAdapter.NewMappingRepository(app.redisClient)
	if err != nil {
		return err
	}

	app.metricsRepo = redisAdapter.NewMetricsRepository(
		app.redisClient,
		app.config.MetricsFlushInterval,
		time.Duration(app.config.MetricsMaxBufferSize),
	)
	return nil
}
