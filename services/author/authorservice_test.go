package services

import (
	"errors"
	"github/brunojoenk/golang-test/models"
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

	authorServiceTest.getAllAuthorsRepository = func(filter models.GetAuthorsFilter) ([]models.Author, error) {
		return []models.Author{{Id: 5, Name: "Joenk"}}, nil
	}

	resp, err := authorServiceTest.GetAllAuthors(models.GetAuthorsFilter{})
	require.NoError(t, err)
	require.Nil(t, deep.Equal([]models.AuthorResponse{{Id: 5, Name: "Joenk"}}, resp.Authors))
	require.Equal(t, resp.Pagination.Limit, 10)
	require.Equal(t, resp.Pagination.Page, 1)
}

func TestGetAllAuthorsError(t *testing.T) {

	authorServiceTest.getAllAuthorsRepository = func(filter models.GetAuthorsFilter) ([]models.Author, error) {
		return nil, errors.New("Error on test")
	}

	_, err := authorServiceTest.GetAllAuthors(models.GetAuthorsFilter{})
	require.Error(t, err)
}

func TestImportAuthorsFromCSVFile(t *testing.T) {
	resp, err := authorServiceTest.ImportAuthorsFromCSVFile("../../data/authorstest.csv")
	require.NoError(t, err)
	require.Equal(t, len(resp), 6)
}

func TestImportAuthorsFromCSVFileError(t *testing.T) {
	_, err := authorServiceTest.ImportAuthorsFromCSVFile("anywhere.csv")
	require.Error(t, err)
}
