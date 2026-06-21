package config

type JWTConfig struct {
	Secret string `env:"SECRET" validate:"required,len=64"`
	Issuer string `env:"ISSUER" validate:"required"`
}
