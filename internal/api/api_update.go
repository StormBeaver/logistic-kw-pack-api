package api

import (
	"context"

	loggerApp "github.com/StormBeaver/logistic-pack-api/internal/logger"
	pb "github.com/StormBeaver/logistic-pack-api/pkg/logistic-pack-api"
	"github.com/opentracing/opentracing-go"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (o *packAPI) UpdatePackV1(
	ctx context.Context,
	req *pb.UpdatePackV1Request,
) (*pb.UpdatePackV1Response, error) {

	span, ctx := opentracing.StartSpanFromContext(ctx, "api.UpdatePack")
	defer span.Finish()

	logger, err := loggerApp.SetLocalLogger(ctx, o.logger)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	if err := req.Validate(); err != nil {

		logger.Error().Err(err).
			Uint64("id", req.GetPackId()).
			Str("name", req.GetName()).
			Str("describe", req.GetDescribe()).Msg("UpdatePackV1 - invalid argument")

		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	logger.Debug().Uint64("id", req.GetPackId()).
		Str("name", req.GetName()).
		Str("describe", req.GetDescribe()).Msg("update arguments")

	updated, err := o.repo.Update(ctx, req.GetPackId(), req.GetName(), req.GetDescribe())

	if err != nil {
		logger.Error().Err(err).Msg("UpdatePackV1 -- failed")

		return nil, status.Error(codes.Internal, err.Error())
	}

	if !updated {
		logger.Debug().Uint64("packId", req.GetPackId()).Msg("pack not found")
		totalPackNotFound.Inc()
		return nil, status.Error(codes.NotFound, "pack not found")
	}

	logger.Debug().Msg("UpdatePackV1 - success")
	totalPackCUDEvents.Inc()

	return &pb.UpdatePackV1Response{
		Found: true,
	}, nil
}
