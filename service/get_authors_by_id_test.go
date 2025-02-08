package service

import (
	pb "github.com/ordarr/authors/v1"
	"github.com/ordarr/data/core"
	"github.com/stretchr/testify/assert"
	_ "github.com/stretchr/testify/suite"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (suite *AuthorTestSuite) TestGetAuthorById() {
	t := suite.T()

	suite.Run("ReturnsPopulatedAuthor", func() {
		suite.mockRepo.(*MockRepo).On("GetByID", []string{"12345"}).Return([]*core.Author{
			{
				BaseTable: core.BaseTable{ID: "12345"},
				Name:      "Name One",
			},
		}, nil)

		out, _ := suite.client.GetAuthors(suite.ctx, &pb.GetAuthorsRequest{Ids: []string{"12345"}})

		assert.NotNil(t, out)
		assert.Equal(t, "Name One", out.Content[0].Name)
	})

	suite.Run("ErrorWhenAuthorDoesntExist", func() {
		suite.mockRepo.(*MockRepo).On("GetByID", []string{"12345"}).Return(nil, status.Error(codes.NotFound, "author not found"))
		_, err := suite.client.GetAuthors(suite.ctx, &pb.GetAuthorsRequest{Ids: []string{"12345"}})

		assert.NotNil(t, err)
		assert.ErrorIs(t, err, status.Error(codes.NotFound, "author not found"))
	})
}
