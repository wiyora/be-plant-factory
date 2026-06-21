package logger

import (
	"context"

	"github.com/rs/zerolog"
)

func WithLayerCtx(ctx context.Context, layer Layer) zerolog.Logger {
	return zerolog.Ctx(ctx).With().Str("layer", string(layer)).Logger()
}

func WithLayer(log *zerolog.Logger, layer Layer) zerolog.Logger {
	return log.With().Str("layer", string(layer)).Logger()
}
