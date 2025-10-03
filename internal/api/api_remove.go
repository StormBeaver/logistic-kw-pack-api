package api

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	loggerApp "github.com/StormBeaver/logistic-pack-api/internal/logger"
	pb "github.com/StormBeaver/logistic-pack-api/pkg/logistic-pack-api"
	"github.com/opentracing/opentracing-go"
)

// RemovePackV1 - remove a pack by ID
func (o *packAPI) RemovePackV1(
	ctx context.Context,
	req *pb.RemovePackV1Request,
) (*pb.RemovePackV1Response, error) {

	span, ctx := opentracing.StartSpanFromContext(ctx, "api.RemovePack")
	defer span.Finish()

	logger, err := loggerApp.SetLocalLogger(ctx, o.logger)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	if err := req.Validate(); err != nil {
		logger.Error().Err(err).Msg("RemovePackV1 - invalid argument")

		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	deleted, err := o.repo.Remove(ctx, req.GetPackId())
	if err != nil {
		logger.Error().Err(err).Msg("RemovePackV1 -- failed")

		return nil, status.Error(codes.Internal, err.Error())
	}

	if !deleted {
		logger.Debug().Uint64("packId", req.GetPackId()).Msg("pack not found")
		totalPackNotFound.Inc()

		return nil, status.Error(codes.NotFound, "pack not found")
	}

	logger.Debug().Msg("RemovePackV1 - success")
	totalPackCUDEvents.Inc()

	return &pb.RemovePackV1Response{
		Found: true,
	}, nil
}
