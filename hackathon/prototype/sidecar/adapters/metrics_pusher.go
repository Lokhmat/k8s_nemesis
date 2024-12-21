package adapters

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"sidecar/domain"
	"sidecar/port"
)

type MetricsPusherAdapter struct {
	apiURL string
}

func NewMetricsPusherAdapter(apiURL string) port.MetricsPusher {
	return &MetricsPusherAdapter{apiURL: apiURL}
}

func (m *MetricsPusherAdapter) Push(metrics *domain.Metrics) error {
	data, err := json.Marshal(metrics)
	if err != nil {
		return err
	}

	resp, err := http.Post(m.apiURL, "application/json", bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to push metrics, status: %d", resp.StatusCode)
	}
	return nil
}
