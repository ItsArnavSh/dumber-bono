package qdrant

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "github.com/qdrant/go-client/qdrant"
)

func NewQdrantClient(ctx context.Context) (pb.CollectionsClient, *grpc.ClientConn, error) {
	conn, err := grpc.NewClient(
		"localhost:6334",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, nil, fmt.Errorf("qdrant connect: %w", err)
	}

	client := pb.NewCollectionsClient(conn)

	if _, err := client.List(ctx, &pb.ListCollectionsRequest{}); err != nil {
		conn.Close()
		return nil, nil, fmt.Errorf("qdrant list collections: %w", err)
	}

	return client, conn, nil
}
