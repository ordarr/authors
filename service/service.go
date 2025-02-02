package service

import (
	"context"
	"github.com/jinzhu/copier"
	pb "github.com/ordarr/authors/v1"
	"github.com/ordarr/data/core"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"
	"google.golang.org/protobuf/types/known/emptypb"
	"log"
	"net"
)

type authorsServer struct {
	pb.UnimplementedAuthorsServer
	repo core.AuthorRepository
}

type AuthorsResult []*pb.Author

func (s *authorsServer) GetAuthors(_ context.Context, _ *emptypb.Empty) (*pb.AuthorsResponse, error) {
	result := &AuthorsResult{}
	if copier.Copy(&result, s.repo.GetAll()) != nil {
		return nil, status.Error(codes.Unknown, "unknown error")
	}
	return &pb.AuthorsResponse{Content: *result}, nil
}

func (s *authorsServer) GetAuthorByName(_ context.Context, request *pb.ValueRequest) (*pb.AuthorResponse, error) {
	result := &pb.Author{}
	author := s.repo.GetByName(request.Value)
	if author.ID == "" || copier.Copy(&result, author) != nil {
		return nil, status.Error(codes.NotFound, "author not found")
	}
	return &pb.AuthorResponse{Content: result}, nil
}

func (s *authorsServer) GetAuthorById(_ context.Context, request *pb.ValueRequest) (*pb.AuthorResponse, error) {
	result := &pb.Author{}
	author := s.repo.GetById(request.Value)
	if author.ID == "" || copier.Copy(&result, author) != nil {
		return nil, status.Error(codes.NotFound, "author not found")
	}
	return &pb.AuthorResponse{Content: result}, nil
}

func NewServer(repository core.AuthorRepository) pb.AuthorsServer {
	t := &authorsServer{
		repo: repository,
	}
	return t
}

func Server(repository core.AuthorRepository) (pb.AuthorsClient, func()) {
	buffer := 101024 * 1024
	lis := bufconn.Listen(buffer)

	baseServer := grpc.NewServer()
	pb.RegisterAuthorsServer(baseServer, NewServer(repository))
	go func() {
		if err := baseServer.Serve(lis); err != nil {
			log.Printf("error serving server: %v", err)
		}
	}()

	conn, err := grpc.NewClient("localhost:8080",
		grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) {
			return lis.Dial()
		}), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("error connecting to server: %v", err)
	}

	closer := func() {
		err := lis.Close()
		if err != nil {
			log.Printf("error closing listener: %v", err)
		}
		baseServer.Stop()
	}

	client := pb.NewAuthorsClient(conn)

	return client, closer
}
