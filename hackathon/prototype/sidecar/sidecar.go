package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
)

type Metrics struct {
	PodName  string  `json:"pod_name"`
	CPUUsage float64 `json:"cpu_usage"`
	RAMUsage float64 `json:"ram_usage"`
}

func getMetrics() (*Metrics, error) {
	cpuPercent, err := cpu.Percent(0, false)
	if err != nil {
		return nil, err
	}

	memStats, err := mem.VirtualMemory()
	if err != nil {
		return nil, err
	}

	return &Metrics{
		PodName:  os.Getenv("HOSTNAME"), // Kubernetes sets this automatically
		CPUUsage: cpuPercent[0],
		RAMUsage: memStats.UsedPercent,
	}, nil
}

func pushMetrics(apiURL string, metrics *Metrics) error {
	data, err := json.Marshal(metrics)
	if err != nil {
		return err
	}

	resp, err := http.Post(apiURL, "application/json", bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to push metrics, status: %d", resp.StatusCode)
	}
	return nil
}

func main() {
	apiURL := os.Getenv("API_URL")
	if apiURL == "" {
		apiURL = "http://localhost:8080/metrics"
	}

	for {
		metrics, err := getMetrics()
		if err != nil {
			fmt.Printf("Error getting metrics: %v\n", err)
			time.Sleep(10 * time.Second)
			continue
		}

		fmt.Printf("Collected Metrics: %+v\n", *metrics)

		err = pushMetrics(apiURL, metrics)
		if err != nil {
			fmt.Printf("Error pushing metrics: %v\n", err)
		}

		time.Sleep(10 * time.Second)
	}
}
