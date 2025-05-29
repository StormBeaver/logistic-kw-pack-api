package api

import (
	"context"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/stormbeaver/logistic-pack-api/pkg/logistic-pack-api"
)

// CreatePackv1 - Create a new pack
func (o *packAPI) CreatePackV1(
	ctx context.Context,
	req *pb.CreatePackV1Request,
) (*pb.CreatePackV1Response, error) {

	if err := req.Validate(); err != nil {
		log.Error().Err(err).Msg("CreatePackV1 - invalid argument")

		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	pack, err := o.repo.CreatePack(ctx, req.GetName())
	if err != nil {
		log.Error().Err(err).Msg("CreatePackV1 -- failed")

		return nil, status.Error(codes.Internal, err.Error())
	}

	if pack == nil {
		log.Debug().Str("packName", req.GetName()).Msg("pack don't create")
		totalPackNotFound.Inc()

		return nil, status.Error(codes.NotFound, "pack don't create")
	}

	log.Debug().Msg("CreatePackV1 - success")

	return &pb.CreatePackV1Response{
		PackId: pack.ID,
	}, nil
}
