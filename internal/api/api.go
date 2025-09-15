package api

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/rs/zerolog"

	"github.com/stormbeaver/logistic-pack-api/internal/repo"

	pb "github.com/stormbeaver/logistic-pack-api/pkg/logistic-pack-api"
)

var (
	totalPackNotFound = promauto.NewCounter(prometheus.CounterOpts{
		Namespace: "logistic_pack_api",
		Name:      "not_found_total",
		Help:      "Total number of packs that were not found",
	})

	totalPackCUDEvents = promauto.NewCounter(prometheus.CounterOpts{
		Namespace: "logistic_pack_api",
		Name:      "cud_found_total",
		Help:      "Total number of create/update/delete events",
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
