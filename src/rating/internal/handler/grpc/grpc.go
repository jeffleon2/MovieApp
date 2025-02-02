package grpc

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"movieexample.com/src/gen"
	"movieexample.com/src/rating/internal/controller/rating"
	"movieexample.com/src/rating/pkg/model"
)

type Handler struct {
	gen.UnimplementedRatingServiceServer
	ctrl *rating.Controller
}

func New(ctrl *rating.Controller) *Handler {
	return &Handler{ctrl: ctrl}
}

func (h *Handler) GetAggregateRating(ctx context.Context, req *gen.GetAggregatedRatingRequest) (*gen.GetAggregatedRatingResponse, error) {
	if req == nil || req.RecordId == "" || req.RecordType == "" {
		return nil, status.Errorf(codes.InvalidArgument, "nil req or empty id")
	}
	v, err := h.ctrl.GetAggregatedRating(ctx, model.RecordID(req.RecordId), model.RecordType(req.RecordType))
	if err != nil && errors.Is(err, rating.ErrNotFound) {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	} else if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}
	return &gen.GetAggregatedRatingResponse{RatingValue: v}, nil
}

func (h *Handler) PutRating(ctx context.Context, req *gen.PutRatingRequest) (*gen.PutRatingResponse, error) {
	if req == nil || req.RecordId == "" || req.RecordType == "" {
		return nil, status.Errorf(codes.InvalidArgument, "nil req or empty user id or record id")
	}
	if err := h.ctrl.PutRating(ctx, model.RecordID(req.RecordId), model.RecordType(req.RecordType), &model.Rating{
		Value: model.RatingValue(req.RatingValue),
	}); err != nil {
		return nil, err
	}
	return &gen.PutRatingResponse{}, nil
}
