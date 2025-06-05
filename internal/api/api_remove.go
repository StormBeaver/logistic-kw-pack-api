package api

import (
	"context"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/stormbeaver/logistic-pack-api/pkg/logistic-pack-api"
)

// RemovePackV1 - remove a pack by ID
func (o *packAPI) RemovePackV1(
	ctx context.Context,
	req *pb.RemovePackV1Request,
) (*pb.RemovePackV1Response, error) {

	if err := req.Validate(); err != nil {
		log.Error().Err(err).Msg("RemovePackV1 - invalid argument")

		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	pack, err := o.repo.RemovePack(ctx, req.GetPackId())
	if err != nil {
		log.Error().Err(err).Msg("RemovePackV1 -- failed")

		return nil, status.Error(codes.Internal, err.Error())
	}

	if pack == nil {
		log.Debug().Uint64("packId", req.GetPackId()).Msg("pack not found")
		totalPackNotFound.Inc()

		return nil, status.Error(codes.NotFound, "pack not found")
	}

	log.Debug().Msg("RemovePackV1 - success")

	return &pb.RemovePackV1Response{
		Found: true,
	}, nil
}
