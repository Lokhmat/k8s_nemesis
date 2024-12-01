package usecase

import (
	"autoscaler/internal/domain/entities"
	"autoscaler/internal/port/repositories"
	"autoscaler/internal/port/services"
	"log"
	"net/http"
	"strings"
	"time"
)

type HandleRequestUseCase struct {
	mappingRepo      repositories.MappingRepository
	deploymentClient services.DeploymentClient
	metricsRepo      repositories.MetricsRepository
	proxy            func(http.ResponseWriter, *http.Request, *entities.Deployment)
}

func NewHandleRequestUseCase(
	mappingRepo repositories.MappingRepository,
	deploymentClient services.DeploymentClient,
	metricsRepo repositories.MetricsRepository,
	proxy func(http.ResponseWriter, *http.Request, *entities.Deployment),
) *HandleRequestUseCase {
	return &HandleRequestUseCase{
		mappingRepo:      mappingRepo,
		deploymentClient: deploymentClient,
		metricsRepo:      metricsRepo,
		proxy:            proxy,
	}
}

func (uc *HandleRequestUseCase) HandleRequest(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	path := r.URL.Path
	var apiPrefix string

	apiPrefix = uc.extractAPIPrefix(path)
	if apiPrefix == "" {
		http.Error(w, "API prefix not recognized", http.StatusNotFound)
		return
	}

	deployment, err := uc.mappingRepo.GetDeploymentByAPIPrefix(ctx, apiPrefix)
	if err != nil {
		log.Printf("Error getting deployment for API prefix %s: %v", apiPrefix, err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	if deployment == nil {
		http.Error(w, "Deployment not found", http.StatusNotFound)
		return
	}

	if err = uc.metricsRepo.UpdateLastRequestTime(ctx, deployment.Name); err != nil {
		log.Printf("Error updating last request time: %v", err)
	}

	available, err := uc.deploymentClient.ArePodsAvailable(ctx, deployment)
	if err != nil {
		log.Printf("Error checking pod availability: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if !available {
		err = uc.deploymentClient.ScaleDeployment(ctx, deployment, 1)
		if err != nil {
			log.Printf("Error scaling deployment %s: %v", deployment.Name, err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		timeout := time.After(2 * time.Minute)
		tick := time.Tick(500 * time.Millisecond)
		for {
			select {
			case <-timeout:
				http.Error(w, "Timeout waiting for pods", http.StatusGatewayTimeout)
				return
			case <-tick:
				available, _ := uc.deploymentClient.ArePodsAvailable(ctx, deployment)
				if available {
					goto Proxy // бу, испугался? не бойся я goto, я тебя не обижу
				}
			}
		}
	}

Proxy:
	uc.proxy(w, r, deployment)
}

func (uc *HandleRequestUseCase) extractAPIPrefix(path string) string {
	segments := strings.Split(strings.TrimPrefix(path, "/"), "/")
	if len(segments) > 0 {
		return "/" + segments[0]
	}
	return ""
}
