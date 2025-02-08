package service

import (
	"context"
	pb "github.com/ordarr/authors/v1"
	"github.com/ordarr/data/core"
	"github.com/stretchr/testify/suite"
	"os"
	"testing"
)

type AuthorTestSuite struct {
	suite.Suite
	ctx      context.Context
	client   pb.AuthorsClient
	mockRepo core.IAuthorRepository
	closer   func()
}

func (suite *AuthorTestSuite) SetupSubTest() {
	suite.ctx = context.Background()
	suite.mockRepo = &MockRepo{}
	suite.client, suite.closer = CreateClient(suite.mockRepo)
}

func (suite *AuthorTestSuite) TearDownSubTest() {
	_ = os.Remove("ordarr.db")
	suite.closer()
}

func TestAuthorTestSuite(t *testing.T) {
	suite.Run(t, new(AuthorTestSuite))
}
