package api

import (
	"context"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/stormbeaver/logistic-pack-api/pkg/logistic-pack-api"
)

// AddPackv1 - Create a new pack
func (o *packAPI) AddPackV1(
	ctx context.Context,
	req *pb.AddPackV1Request,
) (*pb.AddPackV1Response, error) {

	if err := req.Validate(); err != nil {
		log.Error().Err(err).Msg("AddPackV1 - invalid argument")

		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	packId, err := o.repo.Add(ctx, req.GetName())
	if err != nil {
		log.Error().Err(err).Msg("AddPackV1 -- failed")

		return nil, status.Error(codes.Internal, err.Error())
	}

	if packId == 0 {
		log.Debug().Str("packName", req.GetName()).Msg("pack don't create")
		totalPackNotFound.Inc()

		return nil, status.Error(codes.NotFound, "pack don't create")
	}

	log.Debug().Msg("CreatePackV1 - success")

	return &pb.AddPackV1Response{
		PackId: packId,
	}, nil
}
