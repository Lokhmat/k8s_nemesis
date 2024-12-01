package redis

import (
	"autoscaler/internal/port/repositories"
	"context"
	"log"
	"strconv"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
)

type responseTimeEntry struct {
	timestamp    int64
	responseTime int64
}

type metricsRepository struct {
	redisClient       *redis.Client
	lastRequestTimes  map[string]int64
	responseTimeQueue map[string][]responseTimeEntry
	queueMutex        sync.Mutex
	flushInterval     time.Duration
	ctx               context.Context
	cancelFunc        context.CancelFunc
	timeWindow        time.Duration // e.g., 10 minutes
}

func NewMetricsRepository(redisClient *redis.Client, flushInterval, timeWindow time.Duration) repositories.MetricsRepository {
	ctx, cancelFunc := context.WithCancel(context.Background())
	mr := &metricsRepository{
		redisClient:       redisClient,
		lastRequestTimes:  make(map[string]int64),
		responseTimeQueue: make(map[string][]responseTimeEntry),
		flushInterval:     flushInterval,
		ctx:               ctx,
		cancelFunc:        cancelFunc,
		timeWindow:        timeWindow,
	}

	go mr.startFlushLoop()

	return mr
}

func (mr *metricsRepository) UpdateLastRequestTime(_ context.Context, deploymentName string) error {
	now := time.Now().Unix()
	mr.queueMutex.Lock()
	mr.lastRequestTimes[deploymentName] = now
	mr.queueMutex.Unlock()
	return nil
}

func (mr *metricsRepository) UpdateResponseTime(_ context.Context, deploymentName string, responseTime int64) error {
	now := time.Now().Unix()
	entry := responseTimeEntry{
		timestamp:    now,
		responseTime: responseTime,
	}
	mr.queueMutex.Lock()
	mr.responseTimeQueue[deploymentName] = append(mr.responseTimeQueue[deploymentName], entry)
	mr.queueMutex.Unlock()
	return nil
}

func (mr *metricsRepository) startFlushLoop() {
	ticker := time.NewTicker(mr.flushInterval)
	defer ticker.Stop()
	for {
		select {
		case <-mr.ctx.Done():
			mr.flushData()
			return
		case <-ticker.C:
			mr.flushData()
		}
	}
}

func (mr *metricsRepository) flushData() {
	mr.queueMutex.Lock()
	defer mr.queueMutex.Unlock()

	now := time.Now().Unix()
	cutoff := now - int64(mr.timeWindow.Seconds())

	pipeline := mr.redisClient.Pipeline()

	for deploymentName, times := range mr.responseTimeQueue {
		// Remove outdated entries
		validTimes := make([]responseTimeEntry, 0)
		var responseTimeSum int64
		var responseTimeCount int64

		for _, entry := range times {
			if entry.timestamp >= cutoff {
				validTimes = append(validTimes, entry)
				responseTimeSum += entry.responseTime
				responseTimeCount++
			}
		}

		// Update the queue with valid entries
		mr.responseTimeQueue[deploymentName] = validTimes

		// Update Redis with the average response time
		if responseTimeCount > 0 {
			averageResponseTime := responseTimeSum / responseTimeCount
			err := pipeline.HSet(mr.ctx, deploymentName, "averageResponseTime", strconv.FormatInt(averageResponseTime, 10)).Err()
			if err != nil {
				log.Printf("Error updating average response time in Redis: %v", err)
			}
		} else {
			// No valid data, remove the key
			err := pipeline.HDel(mr.ctx, deploymentName, "averageResponseTime").Err()
			if err != nil {
				log.Printf("Error deleting average response time from Redis: %v", err)
			}
		}
	}

	// Update last request times
	for deploymentName, lastRequestTime := range mr.lastRequestTimes {
		err := pipeline.HSet(mr.ctx, deploymentName, "lastRequestTime", strconv.FormatInt(lastRequestTime, 10)).Err()
		if err != nil {
			log.Printf("Error updating last request time in Redis: %v", err)
		}
	}
	mr.lastRequestTimes = make(map[string]int64)

	_, err := pipeline.Exec(mr.ctx)
	if err != nil {
		log.Printf("Error executing pipeline in Redis: %v", err)
	}
}

func (mr *metricsRepository) Shutdown() {
	mr.cancelFunc()
}
