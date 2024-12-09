package models

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	redis "github.com/go-redis/redis/v8"
)

type ScalingRule struct {
	ID         string      `json:"id"`
	Action     Action      `json:"action"`
	Thresholds []Threshold `json:"thresholds"`
}

type Threshold struct {
	MetricType MetricType `json:"metricType"`
	Value      float64    `json:"value"`
	Operator   Operator   `json:"operator"`
}

type Podmetric struct {
	ID                  string    `json:"id"`
	CPUUsage            float64   `json:"cpuUsage"`
	MemoryUsage         float64   `json:"memoryUsage"`
	DiskUsage           float64   `json:"diskUsage"`
	AverageResponseTime float64   `json:"averageResponseTime"`
	ReceivedAt          time.Time `json:"receivedAt"`
}

type Operator string
type MetricType string
type Action string

const (
	OperatorGreaterThan Operator = "GREATER_THAN"
	OperatorLessThan    Operator = "LESS_THAN"
	OperatorEqualTo     Operator = "EQUAL_TO"

	MetricTypeCPUUsage            MetricType = "CPU_USAGE"
	MetricTypeMemoryUsage         MetricType = "MEMORY_USAGE"
	MetricTypeDiskUsage           MetricType = "DISK_USAGE"
	MetricTypeAverageResponseTime MetricType = "AVERAGE_RESPONSE_TIME"

	ActionIncreaseInstances       Action = "INCREASE_INSTANCES"
	ActionDecreaseInstances       Action = "DECREASE_INSTANCES"
	ActionDecreaseInstancesToZero Action = "DECREASE_INSTANCES_TO_ZERO"
)

func saveScalingRule(rdb *redis.Client, scalingRule ScalingRule) error {
	key := fmt.Sprintf("scaling_rule:%s", scalingRule.ID)
	data, err := json.Marshal(scalingRule)
	if err != nil {
		return err
	}
	ctx := context.Background()
	err = rdb.Set(ctx, key, data, 0).Err()
	if err != nil {
		return err
	}

	return nil
}

func getScalingRule(rdb *redis.Client, id string) (*ScalingRule, error) {
	key := fmt.Sprintf("scaling_rule:%s", id)
	ctx := context.Background()
	data, err := rdb.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	var scalingRule ScalingRule
	err = json.Unmarshal([]byte(data), &scalingRule)
	if err != nil {
		return nil, err
	}

	return &scalingRule, nil
}

func getPodmetricsDuration(rdb *redis.Client, id string, duration time.Duration) ([]Podmetric, error) {
	key := fmt.Sprintf("podmetric:%s", id)

	endTime := time.Now().Unix()
	startTime := endTime - int64(duration.Seconds())

	zRangeBy := &redis.ZRangeBy{
		Min: fmt.Sprintf("%d", startTime),
		Max: fmt.Sprintf("%d", endTime),
	}

	ctx := context.Background()
	results, err := rdb.ZRangeByScoreWithScores(ctx, key, zRangeBy).Result()
	if err != nil {
		return nil, err
	}

	var podmetrics []Podmetric
	for _, result := range results {
		var podmetric Podmetric
		err = json.Unmarshal([]byte(result.Member.(string)), &podmetric)
		if err != nil {
			return nil, err
		}
		podmetrics = append(podmetrics, podmetric)
	}

	return podmetrics, nil
}
