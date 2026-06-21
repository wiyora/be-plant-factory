package config

import (
	"fmt"
	"time"
)

type DatabaseConfig struct {
	Host            string        `env:"HOST" validate:"required" default:"localhost"`
	Port            string        `env:"PORT" validate:"required" default:"5432"`
	User            string        `env:"USER" validate:"required"`
	Password        string        `env:"PASSWORD" validate:"required"`
	Name            string        `env:"NAME" validate:"required"`
	SSLMode         string        `env:"SSL_MODE" validate:"required,oneof=disable require verify-ca verify-full" default:"disable"`
	MaxOpenConns    int32         `env:"MAX_OPEN_CONNS" validate:"required,gte=1" default:"25"`
	MinIdleConns    int32         `env:"MIN_IDLE_CONNS" validate:"required,gte=0,ltefield=MaxOpenConns" default:"5"`
	MaxConnLifetime time.Duration `env:"MAX_CONN_LIFETIME" validate:"required,gt=0" default:"30m"`
	MaxConnIdleTime time.Duration `env:"MAX_CONN_IDLE_TIME" validate:"required,gt=0" default:"5m"`
}

func (c DatabaseConfig) DSN() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		c.User,
		c.Password,
		c.Host,
		c.Port,
		c.Name,
		c.SSLMode,
	)
}
