package application

import (
	"autoscaler/config"
	httpAdapter "autoscaler/internal/adapters/http"
	"autoscaler/internal/port/repositories"
	"autoscaler/internal/port/services"
	"context"
	"github.com/go-redis/redis/v8"
	kube "k8s.io/client-go/kubernetes"
	"log"
)

type Application struct {
	config           *config.Config
	redisClient      *redis.Client
	clientSet        *kube.Clientset
	mappingRepo      repositories.MappingRepository
	metricsRepo      repositories.MetricsRepository
	deploymentClient services.DeploymentClient
	server           *httpAdapter.Server
	cancelFuncs      []context.CancelFunc
}

func NewApplication(ctx context.Context) (*Application, error) {
	app := &Application{
		cancelFuncs: make([]context.CancelFunc, 0),
	}

	if err := app.initDependencies(ctx); err != nil {
		return nil, err
	}

	return app, nil
}

func (app *Application) initDependencies(ctx context.Context) error {
	inits := []func(context.Context) error{
		app.initConfig,
		app.initRedisClient,
		app.initKubernetesClient,
		app.initRepositories,
		app.initServices,
		app.initUseCases,
	}

	for _, initFunc := range inits {
		if err := initFunc(ctx); err != nil {
			return err
		}
	}

	return nil
}

func (app *Application) Start() error {
	ctx, cancel := context.WithCancel(context.Background())
	app.cancelFuncs = append(app.cancelFuncs, cancel)

	go app.mappingRepo.SubscribeToUpdates(ctx)

	addr := ":" + app.config.HTTPPort
	log.Printf("Server started on port %s", app.config.HTTPPort)
	return app.server.Start(addr)
}

func (app *Application) Shutdown(ctx context.Context) error {
	for _, cancel := range app.cancelFuncs {
		cancel()
	}

	if err := app.server.Shutdown(ctx); err != nil {
		return err
	}
	app.metricsRepo.Shutdown()

	if err := app.redisClient.Close(); err != nil {
		return err
	}

	log.Println("Application shutdown complete")
	return nil
}

func (app *Application) initConfig(_ context.Context) error {
	conf, err := config.LoadConfig()
	if err != nil {
		return err
	}
	app.config = conf
	return nil
}
