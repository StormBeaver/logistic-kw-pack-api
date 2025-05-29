package api

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"

	"github.com/stormbeaver/logistic-pack-api/internal/repo"

	pb "github.com/stormbeaver/logistic-pack-api/pkg/logistic-pack-api"
)

var (
	totalPackNotFound = promauto.NewCounter(prometheus.CounterOpts{
		Name: "logistic_pack_api_pack_not_found_total",
		Help: "Total number of packs that were not found",
	})
)

type packAPI struct {
	pb.UnimplementedLogisticPackApiServiceServer
	repo repo.Repo
}

// NewPackAPI returns api of logistic-pack-api service
func NewPackAPI(r repo.Repo) pb.LogisticPackApiServiceServer {
	return &packAPI{repo: r}
}
