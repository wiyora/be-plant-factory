package config

type FiberConfig struct {
	BodyLimit      int        `env:"BODY_LIMIT" validate:"required,gt=0" default:"4194304"` // 4 MB
	Prefork        bool       `env:"PREFORK" default:"false"`
	TrustedProxies StringList `env:"TRUSTED_PROXIES" validate:"required,unique" default:"[\"127.0.0.1\", \"::1\"]"`
	CompressLevel  int        `env:"COMPRESS_LEVEL" validate:"required,gte=-1,lte=2" default:"1"`
}

type CORSConfig struct {
	AllowCredentials bool       `env:"ALLOW_CREDENTIALS" default:"false"`
	ExposeHeaders    StringList `env:"EXPOSE_HEADERS" validate:"required,unique" default:"[\"Content-Type\", \"Authorization\", \"X-Request-ID\"]"`
	AllowOrigins     StringList `env:"ALLOW_ORIGINS" validate:"required,unique" default:"[\"*\"]"`
	AllowMethods     StringList `env:"ALLOW_METHODS" validate:"required,unique,dive,oneof=GET POST PUT PATCH DELETE OPTIONS" default:"[\"GET\", \"POST\", \"PUT\", \"PATCH\", \"DELETE\", \"OPTIONS\"]"`
	AllowHeaders     StringList `env:"ALLOW_HEADERS" validate:"required,unique" default:"[\"Origin\", \"Content-Type\", \"Accept\", \"Authorization\", \"X-Request-ID\"]"`
}
