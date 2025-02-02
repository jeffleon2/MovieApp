package grpc

import (
	"context"

	"movieexample.com/src/gen"
	grpcutil "movieexample.com/src/internal"
	"movieexample.com/src/pkg/discovery"
	"movieexample.com/src/rating/pkg/model"
)

// Gateway defines a movie metadata HTTP gateway
type Gateway struct {
	registry discovery.Registry
}

// New creates a new HTTP gateway for a movie metadata service
func New(registry discovery.Registry) *Gateway {
	return &Gateway{registry}
}

func (g *Gateway) GetAggregatedRating(ctx context.Context, recordID model.RecordID, recordType model.RecordType) (float64, error) {
	conn, err := grpcutil.ServiceConnection(ctx, "rating", g.registry)
	if err != nil {
		return 0, err
	}
	defer conn.Close()
	client := gen.NewRatingServiceClient(conn)
	resp, err := client.GetAggregateRating(ctx, &gen.GetAggregatedRatingRequest{RecordId: string(recordID), RecordType: string(recordType)})
	if err != nil {
		return 0, err
	}
	return resp.RatingValue, nil
}

func (g *Gateway) PutRating(ctx context.Context, recordID model.RecordID, recordType model.RecordType, rating *model.Rating) error {
	conn, err := grpcutil.ServiceConnection(ctx, "rating", g.registry)
	if err != nil {
		return err
	}
	defer conn.Close()
	client := gen.NewRatingServiceClient(conn)
	_, err = client.PutRating(ctx, &gen.PutRatingRequest{UserId: string(rating.UserID), RecordId: string(recordID), RecordType: string(recordType), RatingValue: int32(rating.Value)})
	if err != nil {
		return err
	}
	return nil
}
