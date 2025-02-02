package service

import (
	"github.com/google/uuid"
	pb "github.com/ordarr/authors/v1"
	"github.com/stretchr/testify/assert"
	_ "github.com/stretchr/testify/suite"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (suite *AuthorTestSuite) TestGetAuthorByName() {
	t := suite.T()

	suite.Run("ReturnsPopulatedAuthor", func() {
		suite.populate()

		out, _ := suite.client.GetAuthorByName(suite.ctx, &pb.ValueRequest{Value: "Name One"})

		assert.NotNil(t, out)
		assert.NoError(t, uuid.Validate(out.Content.Id))
	})

	suite.Run("ErrorWhenAuthorDoesntExist", func() {
		_, err := suite.client.GetAuthorByName(suite.ctx, &pb.ValueRequest{Value: "some-random-id"})

		assert.NotNil(t, err)
		assert.ErrorIs(t, err, status.Error(codes.NotFound, "author not found"))
	})
}
