package loggerApp

import (
	"context"

	"github.com/opentracing/opentracing-go"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/uber/jaeger-client-go"
	"google.golang.org/grpc/metadata"
)

func LogInit(debugLevel bool) zerolog.Logger {
	if debugLevel {
		return log.Level(zerolog.DebugLevel)
	}
	return log.Level(zerolog.InfoLevel)
}

func SetLocalLogger(ctx context.Context, logger zerolog.Logger) (zerolog.Logger, error) { //TODO: chose a better name
	md, ok := metadata.FromIncomingContext(ctx)
	if ok && (len(md.Get("log-level")) > 0) {
		lvl, err := zerolog.ParseLevel(md.Get("log-level")[0])
		if err != nil {
			logger.Error().Err(err).Msg("can't parse log level")
			return logger, err
		}

		logger.Warn().Msg("change log level to: " + lvl.String())
		return logger.Level(lvl), nil
	}

	if jaegerSpan, ok := opentracing.SpanFromContext(ctx).Context().(jaeger.SpanContext); ok {
		return logger.With().Str("trace-id:", jaegerSpan.TraceID().String()).Logger(), nil
	}

	return logger, nil
}
