package logger

type Layer string

const (
	LayerApp       Layer = "app"
	LayerConfig    Layer = "config"
	LayerValidator Layer = "validator"
	LayerPostgres  Layer = "postgres"
	LayerRedis     Layer = "redis"
	LayerHttp      Layer = "http"

	LayerHTTPServer         Layer = "http_server"
	LayerMiddleware         Layer = "middleware"
	LayerHTTPHandler        Layer = "http_handler"
	LayerUseCase            Layer = "use_case"
	LayerPostgresRepository Layer = "postgres_repository"
	LayerRedisRepository    Layer = "redis_repository"
	LayerCacheRepository    Layer = "cache_repository"
)
