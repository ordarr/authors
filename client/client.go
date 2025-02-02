package client

import (
	"context"
	pb "github.com/ordarr/authors/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/grpclog"
	"log"
	"os"
	"time"
)

func Create(addr *string) (pb.AuthorsClient, context.Context, func()) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	grpclog.SetLoggerV2(grpclog.NewLoggerV2WithVerbosity(os.Stderr, os.Stderr, os.Stderr, 99))

	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to create connection: %v", err)
	}

	connCtx, connCancel := context.WithTimeout(ctx, 20*time.Second)

	return pb.NewAuthorsClient(conn), connCtx, func() {
		connCancel()
		_ = conn.Close()
	}
}
