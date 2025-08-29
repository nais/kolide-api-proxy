package proxy

import (
	"context"

	"github.com/sethvargo/go-envconfig"
)

type Config struct {
	KolideApiToken string `env:"KOLIDE_API_TOKEN,required"`
	ProxyApiToken  string `env:"PROXY_API_TOKEN,required"`
	ListenAddress  string `env:"HTTP_LISTEN_ADDRESS,default=0.0.0.0:8080"`
	LogLevel       string `env:"PROXY_LOG_LEVEL,default=info"`
}

func newConfig(ctx context.Context) (*Config, error) {
	cfg := &Config{}
	if err := envconfig.Process(ctx, cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
