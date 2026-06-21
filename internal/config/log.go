package config

type LogConfig struct {
	Level string `env:"LEVEL" validate:"required,oneof=trace debug info warn error fatal panic" default:"info"`
}
