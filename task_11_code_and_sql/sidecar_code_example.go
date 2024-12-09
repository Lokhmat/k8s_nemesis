package models

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	redis "github.com/go-redis/redis/v8"
)

type Podmetric struct {
	ID                  string    `json:"id"`
	CPUUsage            float64   `json:"cpuUsage"`
	MemoryUsage         float64   `json:"memoryUsage"`
	DiskUsage           float64   `json:"diskUsage"`
	AverageResponseTime float64   `json:"averageResponseTime"`
	ReceivedAt          time.Time `json:"receivedAt"`
}

func savePodmetric(rdb *redis.Client, podmetric Podmetric) error {
	key := fmt.Sprintf("podmetric:%s", podmetric.ID)
	data, err := json.Marshal(podmetric)
	if err != nil {
		return err
	}

	ctx := context.Background()
	score := float64(podmetric.ReceivedAt.Unix())
	err = rdb.ZAdd(ctx, key, &redis.Z{
		Score:  score,
		Member: data,
	}).Err()
	if err != nil {
		return err
	}

	return nil
}

func getPodmetrics(rdb *redis.Client, id string, start, end time.Time) ([]Podmetric, error) {
	key := fmt.Sprintf("podmetric:%s", id)
	zRangeBy := &redis.ZRangeBy{
		Min: fmt.Sprintf("%d", start.Unix()),
		Max: fmt.Sprintf("%d", end.Unix()),
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
