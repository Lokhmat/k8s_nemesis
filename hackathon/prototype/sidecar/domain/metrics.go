package domain

type Metrics struct {
	PodName        string  `json:"pod_name"`
	DeploymentName string  `json:"deployment_name"`
	PodID          string  `json:"pod_id"`
	CPUUsage       float64 `json:"cpu_usage"`
	RAMUsage       float64 `json:"ram_usage"`
}
