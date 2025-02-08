package main

import (
	"flag"
	"fmt"
	"github.com/ordarr/authors/service"
	"github.com/ordarr/data/core"
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
	config, _ := core.BuildConfig()
	s, err := service.Server(&core.AuthorRepository{DB: core.Connect(config)})
	if err != nil {
		log.Fatalf("failed to build config: %v", err)
	}
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
