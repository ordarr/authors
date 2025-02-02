package service

import (
	"context"
	pb "github.com/ordarr/authors/v1"
	"github.com/ordarr/data/core"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
	"os"
	"testing"
)

type AuthorTestSuite struct {
	suite.Suite
	ctx      context.Context
	client   pb.AuthorsClient
	db       *gorm.DB
	populate func() []*core.Author
	closer   func()
}

func insertTestAuthors(db *gorm.DB, ctx context.Context) []*core.Author {
	authorOne := &core.Author{
		Ids: core.Ids{
			Calibre:  1,
			Koreader: 2,
		},
		Name: "Name One",
	}
	authorTwo := &core.Author{
		Ids: core.Ids{
			Calibre:  2,
			Koreader: 3,
		},
		Name: "Name Two",
	}
	session := db.Session(&gorm.Session{Context: ctx})
	session.Create(&authorOne)
	session.Create(&authorTwo)

	return []*core.Author{
		authorOne, authorTwo,
	}
}

func (suite *AuthorTestSuite) SetupSubTest() {
	suite.ctx = context.Background()

	_db := core.Connect(&core.Config{
		Type:    "sqlite",
		Name:    "ordarr.db",
		LogMode: true,
	})

	suite.populate = func() []*core.Author {
		return insertTestAuthors(_db, suite.ctx)
	}

	_client, _closer := Server(core.AuthorRepository{DB: _db})

	suite.db = _db
	suite.client = _client
	suite.closer = _closer
}

func (suite *AuthorTestSuite) TearDownSubTest() {
	_ = os.Remove("ordarr.db")
	suite.closer()
}

func TestAuthorTestSuite(t *testing.T) {
	suite.Run(t, new(AuthorTestSuite))
}
