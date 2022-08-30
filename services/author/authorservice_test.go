package services

import (
	"errors"
	"github/brunojoenk/golang-test/models/dtos"
	"github/brunojoenk/golang-test/models/entities"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-test/deep"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var authorServiceTest *AuthorService

func init() {
	db, _, _ := sqlmock.New()

	dialector := postgres.New(postgres.Config{
		DSN:                  "sqlmock_db_0",
		DriverName:           "postgres",
		Conn:                 db,
		PreferSimpleProtocol: true,
	})
	gormDB, _ := gorm.Open(dialector, &gorm.Config{})

	authorServiceTest = NewAuthorService(gormDB)
}

func TestGetAllAuthors(t *testing.T) {

	authorServiceTest.getAllAuthorsRepository = func(filter dtos.GetAuthorsFilter) ([]entities.Author, error) {
		return []entities.Author{{Id: 5, Name: "Joenk"}}, nil
	}

	resp, err := authorServiceTest.GetAllAuthors(dtos.GetAuthorsFilter{})
	require.NoError(t, err)
	require.Nil(t, deep.Equal([]dtos.AuthorResponse{{Id: 5, Name: "Joenk"}}, resp.Authors))
	require.Equal(t, resp.Pagination.Limit, 10)
	require.Equal(t, resp.Pagination.Page, 1)
}

func TestGetAllAuthorsError(t *testing.T) {

	authorServiceTest.getAllAuthorsRepository = func(filter dtos.GetAuthorsFilter) ([]entities.Author, error) {
		return nil, errors.New("Error on test")
	}

	_, err := authorServiceTest.GetAllAuthors(dtos.GetAuthorsFilter{})
	require.Error(t, err)
}

func TestImportAuthorsFromCSVFile(t *testing.T) {
	authorServiceTest.createAuthorInBatchRepo = func(author []*entities.Author, batchSize int) error {
		return nil
	}
	resp, err := authorServiceTest.ImportAuthorsFromCSVFile("../../data/authorsreduced.csv")
	require.NoError(t, err)
	require.Equal(t, len(resp), 6)
}

func TestImportAuthorsFromCSVFileErronOnCreateAuthorInBatch(t *testing.T) {
	authorServiceTest.createAuthorInBatchRepo = func(author []*entities.Author, batchSize int) error {
		return errors.New("error occurred")
	}
	resp, err := authorServiceTest.ImportAuthorsFromCSVFile("../../data/authorsreduced.csv")
	require.Error(t, err)
	require.Equal(t, len(resp), 0)
}

func TestImportAuthorsFromCSVFileError(t *testing.T) {
	_, err := authorServiceTest.ImportAuthorsFromCSVFile("anyfile")
	require.Error(t, err)
}
