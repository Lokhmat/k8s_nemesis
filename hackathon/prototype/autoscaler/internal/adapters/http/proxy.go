package http

import (
	"autoscaler/internal/domain/entities"
	"autoscaler/internal/port/repositories"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"
)

func ProxyRequest(metricsRepo repositories.MetricsRepository) func(http.ResponseWriter, *http.Request, *entities.Deployment) {
	return func(w http.ResponseWriter, r *http.Request, deployment *entities.Deployment) {
		targetURL, _ := url.Parse("http://" + deployment.Service + ":" + deployment.Port)
		proxy := httputil.NewSingleHostReverseProxy(targetURL)
		r.URL.Path = strings.TrimPrefix(r.URL.Path, deployment.APIPrefix)
		startTime := time.Now()
		proxy.ServeHTTP(w, r)
		duration := time.Since(startTime).Milliseconds()

		ctx := r.Context()
		if err := metricsRepo.UpdateResponseTime(ctx, deployment.Name, duration); err != nil {
			log.Printf("Error updating response time: %v", err)
		}
	}
}
