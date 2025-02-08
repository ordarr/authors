package service

import (
	pb "github.com/ordarr/authors/v1"
	"github.com/ordarr/data/core"
	"github.com/stretchr/testify/assert"
	_ "github.com/stretchr/testify/suite"
)

func (suite *AuthorTestSuite) TestGetAllAuthors() {
	t := suite.T()

	suite.Run("ReturnsPopulatedList", func() {
		suite.mockRepo.(*MockRepo).On("GetAll").Return([]*core.Author{
			{
				BaseTable: core.BaseTable{ID: "12345"},
				Name:      "Author 1",
			},
		}, nil)
		out, _ := suite.client.GetAuthors(suite.ctx, &pb.GetAuthorsRequest{})

		assert.NotNil(t, out)
		assert.Len(t, out.Content, 1)
	})
}
