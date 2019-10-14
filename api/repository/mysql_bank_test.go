package repository_test

import (
	"context"
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/suite"
	"sejastip.id/api"
	"sejastip.id/api/fixture"
	"sejastip.id/api/repository"
)

type mysqlBankTestSuite struct {
	suite.Suite
	mock sqlmock.Sqlmock
	db   *sql.DB

	repo api.BankRepository
}

func (s *mysqlBankTestSuite) SetupSuite() {
	var err error
	s.db, s.mock, err = sqlmock.New()
	if err != nil {
		s.T().Fatalf("error opening mock db: %v", err)
	}

	s.repo = repository.NewMysqlBank(s.db)
}

func (s *mysqlBankTestSuite) TearDownSuite() {
	s.db.Close()
}

func (s *mysqlBankTestSuite) TestCreateBank() {
	bank := fixture.StubbedBank()
	bank.ID = 0

	prep := s.mock.ExpectPrepare("^INSERT INTO banks")
	prep.ExpectExec().WithArgs(
		bank.Name, bank.Image, AnyTime{}, AnyTime{},
	).WillReturnResult(sqlmock.NewResult(1, 1))

	ctx := context.Background()
	err := s.repo.CreateBank(ctx, &bank)

	s.NoError(err)
	s.Equal(bank.ID, int64(1))
}

func TestMysqlBank(t *testing.T) {
	suite.Run(t, new(mysqlBankTestSuite))
}
