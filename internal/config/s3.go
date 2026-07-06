package config

type S3Config struct {
	Endpoint   string `env:"ENDPOINT" validate:"required" default:"http://localhost:9000"`
	AccessKey  string `env:"ACCESS_KEY" validate:"required"`
	SecretKey  string `env:"SECRET_KEY" validate:"required"`
	BucketName string `env:"BUCKET_NAME" validate:"required" default:"plant-factory"`
}
