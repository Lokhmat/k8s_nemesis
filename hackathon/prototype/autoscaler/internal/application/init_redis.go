package application

import (
	"context"
	"github.com/go-redis/redis/v8"
)

func (app *Application) initRedisClient(ctx context.Context) error {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     app.config.RedisAddress,
		Password: app.config.RedisPassword,
		DB:       app.config.RedisDB,
	})

	if _, err := redisClient.Ping(ctx).Result(); err != nil {
		return err
	}

	app.redisClient = redisClient
	return nil
}
