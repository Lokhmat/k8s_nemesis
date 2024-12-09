CREATE TYPE loadbalancer_policy AS ENUM (
    'round_robin',
    'least_connections',
    'random'
);

CREATE TABLE loadbalancer_settings (
    id SERIAL PRIMARY KEY,
    environment VARCHAR(50) NOT NULL,
    policy loadbalancer_policy NOT NULL,
    max_connections INT NOT NULL,
    timeout_seconds INT NOT NULL,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_loadbalancer_settings_env ON loadbalancer_settings(environment);

CREATE TABLE proxy_logs (
    id BIGSERIAL PRIMARY KEY,
    timestamp TIMESTAMP NOT NULL,
    client_ip INET NOT NULL,
    target_url TEXT NOT NULL,
    http_method VARCHAR(10) NOT NULL,
    response_status INT NOT NULL,
    latency_ms INT,
);

CREATE INDEX idx_proxy_logs_timestamp ON proxy_logs(timestamp);
CREATE INDEX idx_proxy_logs_target_url ON proxy_logs(target_url);

