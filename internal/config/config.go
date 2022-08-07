package config

import (
	"github.com/kelseyhightower/envconfig"
)

type MongoConfig struct {
	URI      string `required:"true"`
	User     string `required:"true"`
	Password string `required:"true"`
	Database string `required:"true"`
}

type GRPCConfig struct {
	Port int `required:"true"`
}

type Config struct {
	Mongo MongoConfig `required:"true"`
	GRPC  GRPCConfig  `required:"true"`
}

func New() (*Config, error) {
	cfg := new(Config)

	if err := envconfig.Process("mongo", &cfg.Mongo); err != nil {
		return nil, err
	}

	if err := envconfig.Process("grpc", &cfg.GRPC); err != nil {
		return nil, err
	}

	return cfg, nil
}
