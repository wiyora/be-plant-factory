package logger

type Layer string

const (
	LayerApp       Layer = "app"
	LayerConfig    Layer = "config"
	LayerValidator Layer = "validator"
	LayerPostgres  Layer = "postgres"
	LayerRedis     Layer = "redis"
	LayerS3        Layer = "s3"
	LayerHttp      Layer = "http"
	LayerScheduler Layer = "scheduler"
	LayerMQTT      Layer = "mqtt"
	LayerSocketIO  Layer = "socket_io"

	LayerHTTPServer         Layer = "http_server"
	LayerMiddleware         Layer = "middleware"
	LayerHTTPHandler        Layer = "http_handler"
	LayerUseCase            Layer = "use_case"
	LayerPostgresRepository Layer = "postgres_repository"
	LayerRedisRepository    Layer = "redis_repository"
	LayerCacheRepository    Layer = "cache_repository"
	LayerS3Repository       Layer = "s3_repository"
)

type Section string

const (
	SectionCron Section = "cron"
)
