package api

import (
	"context"

	loggerApp "github.com/StormBeaver/logistic-pack-api/internal/logger"
	"github.com/StormBeaver/logistic-pack-api/internal/model"
	pb "github.com/StormBeaver/logistic-pack-api/pkg/logistic-pack-api"
	"github.com/opentracing/opentracing-go"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// ListPackv1 - list all packs
func (o *packAPI) ListPackV1(
	ctx context.Context,
	req *pb.ListPackV1Request,
) (*pb.ListPackV1Response, error) {

	span, ctx := opentracing.StartSpanFromContext(ctx, "api.ListPacks")
	defer span.Finish()

	logger, err := loggerApp.SetLocalLogger(ctx, o.logger)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	if err := req.Validate(); err != nil {
		logger.Error().Err(err).
			Uint64("cursor", req.GetCursor()).
			Uint64("limit", req.GetLimit()).Msg("UpdatePackV1 - invalid arguments")

		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	logger.Debug().
		Uint64("cursor", req.GetCursor()).
		Uint64("limit", req.GetLimit()).Msg("list arguments")

	packs, err := o.repo.List(ctx, req.GetCursor(), req.GetLimit())
	if err != nil {
		logger.Error().Err(err).Msg("ListPackV1 -- failed")

		return nil, status.Error(codes.Internal, err.Error())
	}

	if packs == nil {
		logger.Debug().Str("packs", "list").Msg("packs not found")
		totalPackNotFound.Inc()

		return nil, status.Error(codes.NotFound, "packs not found")
	}

	logger.Debug().Msg("ListPackV1 - success")

	return &pb.ListPackV1Response{
		Items: convertPacksList(packs),
	}, nil
}

func convertPacksList(packs []*model.Pack) []*pb.Pack {
	result := make([]*pb.Pack, 0, len(packs))
	for _, v := range packs {
		result = append(result, &pb.Pack{
			Id:      v.ID,
			Name:    v.Name,
			Created: timestamppb.New(v.Created),
		})
	}
	return result
}
