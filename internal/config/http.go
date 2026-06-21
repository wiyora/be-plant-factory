package config

import (
	"fmt"
	"time"
)

type HTTPConfig struct {
	Host            string        `env:"HOST" validate:"required" default:"0.0.0.0"`
	Port            string        `env:"PORT" validate:"required" default:"8080"`
	ReadTimeout     time.Duration `env:"READ_TIMEOUT" validate:"required,gt=0" default:"10s"`
	WriteTimeout    time.Duration `env:"WRITE_TIMEOUT" validate:"required,gt=0" default:"30s"`
	IdleTimeout     time.Duration `env:"IDLE_TIMEOUT" validate:"required,gt=0" default:"60s"`
	ShutdownTimeout time.Duration `env:"SHUTDOWN_TIMEOUT" validate:"required,gt=0" default:"15s"`
}

func (c HTTPConfig) Address() string {
	return fmt.Sprintf("%s:%s", c.Host, c.Port)
}

func (c HTTPConfig) AddressWithScheme() string {
	scheme := "http"
	if c.Port == "443" {
		scheme = "https"
	}

	return fmt.Sprintf("%s://%s:%s", scheme, c.Host, c.Port)
}
