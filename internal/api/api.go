package api

import (
	"context"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/rs/zerolog"
	"google.golang.org/grpc/metadata"

	"github.com/stormbeaver/logistic-pack-api/internal/repo"

	pb "github.com/stormbeaver/logistic-pack-api/pkg/logistic-pack-api"
)

var (
	totalPackNotFound = promauto.NewCounter(prometheus.CounterOpts{
		Name: "logistic_pack_api_pack_not_found_total",
		Help: "Total number of packs that were not found",
	})

	totalPackCUDEvents = promauto.NewCounter(prometheus.CounterOpts{
		Name: "logistic_pack_api_cud_fount_total",
		Help: "Total number of create/update/delete events",
	})
)

type packAPI struct {
	pb.UnimplementedLogisticPackApiServiceServer
	repo   repo.Repo
	logger zerolog.Logger
}

// NewPackAPI returns api of logistic-pack-api service
func NewPackAPI(r repo.Repo, l zerolog.Logger) pb.LogisticPackApiServiceServer {
	return &packAPI{repo: r, logger: l}
}

func checkLogLevel(ctx context.Context, logger zerolog.Logger) (zerolog.Logger, error) {
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
	return logger, nil
}
