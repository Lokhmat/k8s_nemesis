package redis

import (
	"autoscaler/internal/domain/entities"
	"autoscaler/internal/port/repositories"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"

	"github.com/go-redis/redis/v8"
)

type mappingRepository struct {
	redisClient           *redis.Client
	apiPrefixToDeployment map[string]*entities.Deployment
	mappingsMutex         sync.RWMutex
}

func NewMappingRepository(redisClient *redis.Client) (repositories.MappingRepository, error) {
	mr := &mappingRepository{
		redisClient:           redisClient,
		apiPrefixToDeployment: make(map[string]*entities.Deployment),
	}
	if err := mr.LoadMappings(context.Background()); err != nil {
		return nil, fmt.Errorf("load mappings: %w", err)
	}
	go mr.SubscribeToUpdates(context.Background())
	return mr, nil
}

func (mr *mappingRepository) GetDeploymentByAPIPrefix(_ context.Context, apiPrefix string) (*entities.Deployment, error) {
	mr.mappingsMutex.RLock()
	defer mr.mappingsMutex.RUnlock()
	deployment, exists := mr.apiPrefixToDeployment[apiPrefix]
	if !exists {
		return nil, nil
	}
	return deployment, nil
}

func (mr *mappingRepository) LoadMappings(ctx context.Context) error {
	data, err := mr.redisClient.HGetAll(ctx, "deployments_info").Result()
	if err != nil {
		return err
	}

	var targetDeployments []entities.Deployment
	mr.mappingsMutex.Lock()
	defer mr.mappingsMutex.Unlock()
	for _, deploymentData := range data {
		err = json.Unmarshal([]byte(deploymentData), &targetDeployments)
		if err != nil {
			log.Printf("Error unmarshalling data from Redis: %v", err)
			continue
		}

		for _, td := range targetDeployments {
			dep := td // create a copy
			mr.apiPrefixToDeployment[td.APIPrefix] = &dep
		}
	}
	return nil
}

func (mr *mappingRepository) SubscribeToUpdates(ctx context.Context) {
	pubsub := mr.redisClient.Subscribe(ctx, "deployments_channel")
	defer pubsub.Close()

	for {
		msg, err := pubsub.ReceiveMessage(ctx)
		if err != nil {
			log.Printf("Error receiving message from Redis Pub/Sub: %v", err)
			continue
		}

		if msg.Channel == "deployments_channel" {
			if msg.Payload == "updated" || msg.Payload == "deleted" {
				mr.LoadMappings(ctx)
			}
		}
	}
}
