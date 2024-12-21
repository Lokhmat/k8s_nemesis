package main

import (
	"os"
	"sidecar/adapters"
	"sidecar/application"
	"time"
)

func main() {
	apiURL := os.Getenv("API_URL")
	if apiURL == "" {
		apiURL = "http://localhost:8080/metrics"
	}

	collector := adapters.NewMetricsCollectorAdapter()
	pusher := adapters.NewMetricsPusherAdapter(apiURL)

	sidecar := application.NewSidecar(collector, pusher, 10*time.Second)
	sidecar.Run()
}
