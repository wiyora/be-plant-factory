package config

import "time"

type AuthConfig struct {
	Google               AuthGoogleConfig `env:"GOOGLE"`
	SocialStateDuration  time.Duration    `env:"SOCIAL_STATE_DURATION" validate:"required,gt=0" default:"10m"`
	AccessTokenDuration  time.Duration    `env:"ACCESS_TOKEN_DURATION" validate:"required,gt=0" default:"15m"`
	RefreshTokenDuration time.Duration    `env:"REFRESH_TOKEN_DURATION" validate:"required,gt=0" default:"168h"`
	CacheUserDuration    time.Duration    `env:"CACHE_USER_DURATION" validate:"required,gt=0" default:"5m"`
	ToleranceDuration    time.Duration    `env:"TOLERANCE_DURATION" validate:"required,gt=0,ltfield=SocialStateDuration" default:"3m"`
}

type AuthGoogleConfig struct {
	ClientID     string `env:"CLIENT_ID" validate:"required"`
	ClientSecret string `env:"CLIENT_SECRET" validate:"required"`
	RedirectURL  string `env:"REDIRECT_URL" validate:"required,url"`
}
