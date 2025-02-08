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
	"log"
	"net"
)

type AuthorsServer struct {
	pb.UnimplementedAuthorsServer
	repo core.IAuthorRepository
}

type AuthorsResult []*pb.Author

func (s *AuthorsServer) GetAuthors(_ context.Context, request *pb.GetAuthorsRequest) (*pb.AuthorsResponse, error) {
	result := &AuthorsResult{}
	var authors []*core.Author
	var err error

	if len(request.Ids) == 0 && len(request.Names) == 0 {
		authors, err = s.repo.GetAll()
	} else if request.Ids != nil {
		authors, err = s.repo.GetByID(request.Ids)
	} else {
		authors, err = s.repo.GetByName(request.Names)
	}
	if err != nil {
		return nil, err
	}

	if copier.Copy(&result, authors) != nil {
		return nil, status.Error(codes.Unknown, "unknown error")
	}
	return &pb.AuthorsResponse{Content: *result}, nil
}

func (s *AuthorsServer) CreateAuthor(_ context.Context, request *pb.CreateAuthorRequest) (*pb.AuthorResponse, error) {
	result := &pb.Author{}
	created := &core.Author{
		Name: request.Name,
	}
	created, err := s.repo.Create(created)
	if err != nil {
		return nil, err
	}
	if copier.Copy(&result, created) != nil {
		return nil, status.Error(codes.Unknown, "unknown error")
	}
	return &pb.AuthorResponse{Content: result}, nil
}

func Server(repository core.IAuthorRepository) (*grpc.Server, error) {
	baseServer := grpc.NewServer()
	pb.RegisterAuthorsServer(baseServer, &AuthorsServer{repo: repository})
	return baseServer, nil
}

func CreateClient(repository core.IAuthorRepository) (pb.AuthorsClient, func()) {
	buffer := 101024 * 1024
	lis := bufconn.Listen(buffer)

	baseServer := grpc.NewServer()

	pb.RegisterAuthorsServer(baseServer, &AuthorsServer{repo: repository})

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
