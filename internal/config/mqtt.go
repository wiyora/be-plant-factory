package config

import "fmt"

type MQTTConfig struct {
	Broker   string `env:"BROKER" validate:"required" default:"tcp://localhost:1883"`
	ClientID string `env:"CLIENT_ID" validate:"required" default:"plant-factory"`
	Username string `env:"USERNAME"`
	Password string `env:"PASSWORD"`
}

func (c MQTTConfig) Address() string {
	return c.Broker
}

func (c MQTTConfig) ClientIdentifier() string {
	if c.ClientID != "" {
		return c.ClientID
	}

	return fmt.Sprintf("plant-factory-%d", 0)
}
