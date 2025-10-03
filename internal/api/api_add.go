package api

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	loggerApp "github.com/StormBeaver/logistic-pack-api/internal/logger"
	pb "github.com/StormBeaver/logistic-pack-api/pkg/logistic-pack-api"
	"github.com/opentracing/opentracing-go"
)

// AddPackv1 - Create a new pack
func (o *packAPI) AddPackV1(
	ctx context.Context,
	req *pb.AddPackV1Request,
) (*pb.AddPackV1Response, error) {

	span, ctx := opentracing.StartSpanFromContext(ctx, "api.AddPack")
	defer span.Finish()

	logger, err := loggerApp.SetLocalLogger(ctx, o.logger)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	if err := req.Validate(); err != nil {
		logger.Error().Err(err).Msg("AddPackV1 - invalid argument")

		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	packId, err := o.repo.Add(ctx, req.GetName())
	if err != nil {
		logger.Error().Err(err).Msg("AddPackV1 -- failed")

		return nil, status.Error(codes.Internal, err.Error())
	}

	if packId == 0 {
		logger.Debug().Str("packName", req.GetName()).Msg("don't create")

		return nil, status.Error(codes.NotFound, "pack don't create")
	}

	logger.Debug().Msg("CreatePackV1 - success")
	totalPackCUDEvents.Inc()

	return &pb.AddPackV1Response{
		PackId: packId,
	}, nil
}
