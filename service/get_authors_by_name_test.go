package service

import (
	pb "github.com/ordarr/authors/v1"
	"github.com/ordarr/data/core"
	"github.com/stretchr/testify/assert"
	_ "github.com/stretchr/testify/suite"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (suite *AuthorTestSuite) TestGetAuthorByName() {
	t := suite.T()

	suite.Run("ReturnsPopulatedAuthor", func() {
		suite.mockRepo.(*MockRepo).On("GetByName", []string{"Name One"}).Return([]*core.Author{
			{
				BaseTable: core.BaseTable{ID: "12345"},
				Name:      "Name One",
			},
		}, nil)

		out, _ := suite.client.GetAuthors(suite.ctx, &pb.GetAuthorsRequest{Names: []string{"Name One"}})

		assert.NotNil(t, out)
		assert.Equal(t, "12345", out.Content[0].Id)
	})

	suite.Run("ErrorWhenAuthorDoesntExist", func() {
		suite.mockRepo.(*MockRepo).On("GetByName", []string{"Name One"}).Return(nil, status.Error(codes.NotFound, "author not found"))
		_, err := suite.client.GetAuthors(suite.ctx, &pb.GetAuthorsRequest{Names: []string{"Name One"}})

		assert.NotNil(t, err)
		assert.ErrorIs(t, err, status.Error(codes.NotFound, "author not found"))
	})
}
