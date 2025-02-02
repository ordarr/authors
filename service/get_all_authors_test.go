package service

import (
	"github.com/stretchr/testify/assert"
	_ "github.com/stretchr/testify/suite"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (suite *AuthorTestSuite) TestGetAllAuthors() {
	t := suite.T()

	suite.Run("ReturnsPopulatedList", func() {
		suite.populate()

		out, _ := suite.client.GetAuthors(suite.ctx, &emptypb.Empty{})

		assert.NotNil(t, out)
		assert.Len(t, out.Content, 2)
	})

	suite.Run("ReturnsEmptyList", func() {
		out, _ := suite.client.GetAuthors(suite.ctx, &emptypb.Empty{})

		assert.NotNil(t, out)
		assert.Len(t, out.Content, 0)
	})
}
