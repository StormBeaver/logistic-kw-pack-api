package api

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	loggerApp "github.com/StormBeaver/logistic-pack-api/internal/logger"
	pb "github.com/StormBeaver/logistic-pack-api/pkg/logistic-pack-api"
	"github.com/opentracing/opentracing-go"
)

// GetPackV1 - Get a pack by ID
func (o *packAPI) GetPackV1(
	ctx context.Context,
	req *pb.GetPackV1Request,
) (*pb.GetPackV1Response, error) {

	span, ctx := opentracing.StartSpanFromContext(ctx, "api.GetPack")
	defer span.Finish()

	logger, err := loggerApp.SetLocalLogger(ctx, o.logger)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	if err := req.Validate(); err != nil {
		logger.Error().Err(err).
			Uint64("id", req.GetPackId()).Msg("GetPackV1 - invalid argument")

		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	logger.Debug().Uint64("id", req.GetPackId()).Msg("get argument")

	pack, err := o.repo.Get(ctx, req.GetPackId())
	if err != nil {
		logger.Error().Err(err).Msg("GetPackV1 -- failed")

		return nil, status.Error(codes.Internal, err.Error())
	}

	if pack == nil {
		logger.Debug().Uint64("packId", req.PackId).Msg("pack not found")
		totalPackNotFound.Inc()

		return nil, status.Error(codes.NotFound, "pack not found")
	}

	logger.Debug().Msg("DescribePackV1 - success")

	return &pb.GetPackV1Response{
		Value: &pb.Pack{
			Id:      pack.ID,
			Name:    pack.Name,
			Created: timestamppb.New(pack.Created),
		},
	}, nil
}
