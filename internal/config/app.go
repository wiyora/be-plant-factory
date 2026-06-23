package config

import (
	"os"
	"strings"
)

type AppConfig struct {
	Env  AppEnv `env:"ENV" validate:"required,oneof=development staging production" default:"development"`
	Name string `env:"NAME" validate:"required" default:"Plant Factory"`
}

type AppEnv string

func (env AppEnv) IsProduction() bool {
	return env == "production"
}

func (env AppEnv) IsDevelopment() bool {
	return env == "development"
}

func (env AppEnv) IsStaging() bool {
	return env == "staging"
}

func (env AppEnv) IsServerEnv() bool {
	return env.IsProduction() || env.IsStaging()
}

func (env AppEnv) String() string {
	return string(env)
}

func MustLoadAppEnv() AppEnv {
	rawEnv := os.Getenv("APP_ENV")
	appEnv := strings.ToLower(strings.TrimSpace(rawEnv))

	switch appEnv {
	case "production":
		return "production"
	case "staging":
		return "staging"
	default:
		return "development"
	}
}
