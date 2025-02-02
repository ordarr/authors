package main

import (
	"flag"
	"fmt"
	"github.com/ordarr/authors/service"
	pb "github.com/ordarr/authors/v1"
	"github.com/ordarr/data/core"
	"google.golang.org/grpc"
	"log"
	"net"
)

var (
	port = flag.Int("port", 50051, "The server port")
)

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	config, err := core.BuildConfig()
	if err != nil {
		log.Fatalf("failed to build config: %v", err)
	}
	connect := core.Connect(config)
	pb.RegisterAuthorsServer(s, service.NewServer(core.AuthorRepository{DB: connect}))
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
