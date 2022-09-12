package repository

import (
	"context"
	"database/sql"
	"github/brunojoenk/golang-test/models/dtos"
	"github/brunojoenk/golang-test/models/entities"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-test/deep"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Suite struct {
	suite.Suite
	DB   *gorm.DB
	mock sqlmock.Sqlmock

	repository *AuthorRepository
}

func (s *Suite) AfterTest(_, _ string) {
	require.NoError(s.T(), s.mock.ExpectationsWereMet())
}

func TestInit(t *testing.T) {
	suite.Run(t, new(Suite))
}

func (s *Suite) SetupSuite() {
	var (
		db  *sql.DB
		err error
	)

	db, s.mock, err = sqlmock.New()
	require.NoError(s.T(), err)

	dialector := postgres.New(postgres.Config{
		DSN:                  "sqlmock_db_0",
		DriverName:           "postgres",
		Conn:                 db,
		PreferSimpleProtocol: true,
	})
	s.DB, err = gorm.Open(dialector, &gorm.Config{})
	require.NoError(s.T(), err)

	s.repository = &AuthorRepository{db: s.DB}
}

func (s *Suite) Test_repository_Get_Author() {
	var (
		id   = 1
		name = "test-name"
	)

	s.mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "authors" WHERE "authors"."id" = $1`)).
		WithArgs(id).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).
			AddRow(id, name))

	res, err := s.repository.GetAuthor(id)

	require.NoError(s.T(), err)
	require.Nil(s.T(), deep.Equal(&entities.Author{Id: id, Name: name}, res))
}

func (s *Suite) Test_repository_Get_Author_Error() {

	s.mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "authors" WHERE "authors"."id" = $1`)).
		WillReturnError(context.Canceled)

	_, err := s.repository.GetAuthor(1)

	require.Error(s.T(), err)
}

func (s *Suite) Test_repository_Get_All_Authors() {
	var (
		id   = 1
		name = "test-name"
	)

	s.mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "authors"`)).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).
			AddRow(id, name))

	res, err := s.repository.GetAllAuthors(dtos.GetAuthorsFilter{})

	require.NoError(s.T(), err)
	require.Nil(s.T(), deep.Equal([]entities.Author{{Id: id, Name: name}}, res))
}

func (s *Suite) Test_repository_Get_All_Authors_Error() {

	s.mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "authors"`)).
		WillReturnError(context.Canceled)

	_, err := s.repository.GetAllAuthors(dtos.GetAuthorsFilter{})

	require.Error(s.T(), err)

}

func (s *Suite) Test_repository_Get_All_Authors_Filter_Name() {
	var (
		id   = 1
		name = "test-name"
	)

	s.mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "authors" WHERE LOWER(name) LIKE $1`)).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).
			AddRow(id, name))

	res, err := s.repository.GetAllAuthors(dtos.GetAuthorsFilter{Name: name})

	require.NoError(s.T(), err)
	require.Nil(s.T(), deep.Equal([]entities.Author{{Id: id, Name: name}}, res))
}

func (s *Suite) Test_repository_Create_Author() {
	var (
		id   = 1
		name = "test-name"
	)

	s.mock.ExpectBegin()

	s.mock.ExpectQuery(regexp.QuoteMeta(
		`INSERT INTO "authors" ("name") VALUES ($1)`)).
		WithArgs(name).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).
			AddRow(id))
	s.mock.ExpectCommit()

	err := s.repository.CreateAuthorInBatch([]entities.Author{{Name: name}}, 1)

	require.NoError(s.T(), err)
}

func (s *Suite) Test_repository_Create_Author_Error() {
	var (
		name = "test-name"
	)

	s.mock.ExpectBegin()

	s.mock.ExpectQuery(regexp.QuoteMeta(
		`INSERT INTO "authors" ("name") VALUES ($1)`)).
		WithArgs(name).
		WillReturnError(context.Canceled)

	s.mock.ExpectRollback()

	err := s.repository.CreateAuthorInBatch([]entities.Author{{Name: name}}, 1)

	require.Error(s.T(), err)
}
