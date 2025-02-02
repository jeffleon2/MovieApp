package grpc

import (
	"context"

	"movieexample.com/src/gen"
	grpcutil "movieexample.com/src/internal"
	"movieexample.com/src/metadata/pkg/model"
	"movieexample.com/src/pkg/discovery"
)

// Gateway defines a movie metadata GRPC gateway
type Gateway struct {
	registry discovery.Registry
}

// New creates a new GRPC gateway for a movie metadata service
func New(registry discovery.Registry) *Gateway {
	return &Gateway{registry}
}

func (g *Gateway) Get(ctx context.Context, id string) (*model.Metadata, error) {
	conn, err := grpcutil.ServiceConnection(ctx, "metadata", g.registry)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	client := gen.NewMetadataServiceClient(conn)
	resp, err := client.GetMetadata(ctx, &gen.GetMetadataRequest{MovieId: id})
	if err != nil {
		return nil, err
	}
	return model.MetadataFromProto(resp.Metadata), nil
}
