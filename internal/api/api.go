package api

import (
	"context"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"route255/logistic-kw-pack-api/internal/repo"

	pb "github.com/ozonmp/omp-pack-api/pkg/omp-pack-api"
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

func (o *packAPI) DescribePackV1(
	ctx context.Context,
	req *pb.DescribePackV1Request,
) (*pb.DescribePackV1Response, error) {

	if err := req.Validate(); err != nil {
		log.Error().Err(err).Msg("DescribePackV1 - invalid argument")

		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	pack, err := o.repo.DescribePack(ctx, req.PackId)
	if err != nil {
		log.Error().Err(err).Msg("DescribePackV1 -- failed")

		return nil, status.Error(codes.Internal, err.Error())
	}

	if pack == nil {
		log.Debug().Uint64("packId", req.PackId).Msg("pack not found")
		totalPackNotFound.Inc()

		return nil, status.Error(codes.NotFound, "pack not found")
	}

	log.Debug().Msg("DescribePackV1 - success")

	return &pb.DescribePackV1Response{
		Value: &pb.Pack{
			Id:  pack.ID,
			Foo: pack.Foo,
		},
	}, nil
}
