package service

import (
	pb "github.com/ordarr/authors/v1"
	"github.com/stretchr/testify/assert"
	_ "github.com/stretchr/testify/suite"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (suite *AuthorTestSuite) TestGetAuthorById() {
	t := suite.T()

	suite.Run("ReturnsPopulatedAuthor", func() {
		inserted := suite.populate()

		out, _ := suite.client.GetAuthorById(suite.ctx, &pb.ValueRequest{Value: inserted[0].ID})

		assert.NotNil(t, out)
		assert.Equal(t, inserted[0].Name, out.Content.Name)
	})

	suite.Run("ErrorWhenAuthorDoesntExist", func() {
		t := suite.T()

		_, err := suite.client.GetAuthorById(suite.ctx, &pb.ValueRequest{Value: "4783e133-d856-43f4-8d38-9e50c5996cad"})

		assert.NotNil(t, err)
		assert.ErrorIs(t, err, status.Error(codes.NotFound, "author not found"))
	})
}
