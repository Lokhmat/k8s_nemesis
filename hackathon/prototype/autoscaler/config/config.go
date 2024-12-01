package config

import (
	"github.com/caarlos0/env/v6"
	"time"
)

type Config struct {
	// Kubernetes configuration
	KubeConfig string        `env:"KUBECONFIG"`
	Namespace  string        `env:"NAMESPACE" envDefault:"default"`
	IdleTime   time.Duration `env:"IDLE_TIME" envDefault:"60s"`

	// Redis configuration
	RedisAddress  string `env:"REDIS_ADDRESS" envDefault:"redis:6379"`
	RedisPassword string `env:"REDIS_PASSWORD" envDefault:""`
	RedisDB       int    `env:"REDIS_DB" envDefault:"0"`

	// HTTP Server configuration
	HTTPPort string `env:"HTTP_PORT" envDefault:"8080"`

	// Metrics buffer configuration
	MetricsFlushInterval time.Duration `env:"METRICS_FLUSH_INTERVAL" envDefault:"5s"`
	MetricsMaxBufferSize int           `env:"METRICS_MAX_BUFFER_SIZE" envDefault:"100"`
}

func LoadConfig() (*Config, error) {
	cfg := &Config{}
	err := env.Parse(cfg)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}
