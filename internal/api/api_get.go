package api

import (
	"context"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	pb "github.com/stormbeaver/logistic-pack-api/pkg/logistic-pack-api"
)

// GetPackV1 - Get a pack by ID
func (o *packAPI) GetPackV1(
	ctx context.Context,
	req *pb.GetPackV1Request,
) (*pb.GetPackV1Response, error) {

	if err := req.Validate(); err != nil {
		log.Error().Err(err).Msg("GetPackV1 - invalid argument")

		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	pack, err := o.repo.Get(ctx, req.GetPackId())
	if err != nil {
		log.Error().Err(err).Msg("GetPackV1 -- failed")

		return nil, status.Error(codes.Internal, err.Error())
	}

	if pack == nil {
		log.Debug().Uint64("packId", req.PackId).Msg("pack not found")
		totalPackNotFound.Inc()

		return nil, status.Error(codes.NotFound, "pack not found")
	}

	log.Debug().Msg("DescribePackV1 - success")

	return &pb.GetPackV1Response{
		Value: &pb.Pack{
			Id:      pack.ID,
			Name:    pack.Name,
			Created: timestamppb.New(pack.Created),
		},
	}, nil
}
