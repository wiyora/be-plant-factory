package config

import "fmt"

type RedisConfig struct {
	Host     string `env:"HOST" validate:"required" default:"localhost"`
	Port     string `env:"PORT" validate:"required" default:"6379"`
	Password string `env:"PASSWORD"`
	DB       int    `env:"DB" default:"0"`
}

func (c RedisConfig) Address() string {
	return fmt.Sprintf("%s:%s", c.Host, c.Port)
}
