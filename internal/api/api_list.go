package api

import (
	"context"

	"github.com/stormbeaver/logistic-pack-api/internal/model"
	pb "github.com/stormbeaver/logistic-pack-api/pkg/logistic-pack-api"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// ListPackv1 - list all packs
func (o *packAPI) ListPackV1(
	ctx context.Context,
	req *pb.ListPackV1Request,
) (*pb.ListPackV1Response, error) {

	logger, err := checkLogLevel(ctx, o.logger)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	if err := req.Validate(); err != nil {
		logger.Error().Err(err).Msg("ListPackV1 - invalid arguments")

		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	pack, err := o.repo.List(ctx, req.GetCursor(), req.GetLimit())
	if err != nil {
		logger.Error().Err(err).Msg("ListPackV1 -- failed")

		return nil, status.Error(codes.Internal, err.Error())
	}

	if pack == nil {
		logger.Debug().Str("packs", "list").Msg("packs not found")
		totalPackNotFound.Inc()

		return nil, status.Error(codes.NotFound, "packs not found")
	}

	logger.Debug().Msg("ListPackV1 - success")

	return &pb.ListPackV1Response{
		Items: convertPackList(pack),
	}, nil
}

func convertPackList(pack []*model.Pack) []*pb.Pack {
	result := make([]*pb.Pack, 0, len(pack))
	for _, v := range pack {
		result = append(result, &pb.Pack{
			Id:      v.ID,
			Name:    v.Name,
			Created: timestamppb.New(v.Created),
		})
	}
	return result
}
