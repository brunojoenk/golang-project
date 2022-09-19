package services

import (
	"errors"
	"github/brunojoenk/golang-test/models/dtos"
	"github/brunojoenk/golang-test/models/entities"
	"testing"

	"github.com/go-test/deep"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	authorrepo "github/brunojoenk/golang-test/repository/author"
)

func TestGetAllAuthors(t *testing.T) {

	authorDbMock := new(authorrepo.AuthorRepositoryMock)

	filter := dtos.GetAuthorsFilter{Pagination: dtos.Pagination{Page: 1, Limit: 10}}
	authorDbMock.On("GetAllAuthors", filter).Return([]entities.Author{{Id: 5, Name: "Joenk"}}, nil)

	authorServiceTest := AuthorService{authorDb: authorDbMock}

	resp, err := authorServiceTest.GetAllAuthors(filter)
	require.NoError(t, err)
	require.Nil(t, deep.Equal([]dtos.AuthorResponse{{Id: 5, Name: "Joenk"}}, resp.Authors))
	require.Equal(t, resp.Pagination.Limit, 10)
	require.Equal(t, resp.Pagination.Page, 1)
}

func TestGetAllAuthorsError(t *testing.T) {

	authorDbMock := new(authorrepo.AuthorRepositoryMock)

	filter := dtos.GetAuthorsFilter{Pagination: dtos.Pagination{Page: 1, Limit: 10}}
	authorDbMock.On("GetAllAuthors", filter).Return([]entities.Author{}, errors.New("Error on test"))

	authorServiceTest := AuthorService{authorDb: authorDbMock}

	_, err := authorServiceTest.GetAllAuthors(filter)
	require.Error(t, err)
}

func TestImportAuthorsFromCSVFile(t *testing.T) {
	authorDbMock := new(authorrepo.AuthorRepositoryMock)

	authorDbMock.On("CreateAuthorInBatch", mock.Anything, 6).Return(nil)

	authorServiceTest := AuthorService{authorDb: authorDbMock}

	resp, err := authorServiceTest.ImportAuthorsFromCSVFile("../../data/authorsreduced.csv")
	require.NoError(t, err)
	require.Equal(t, 6, resp)
}

func TestImportAuthorsFromCSVFileErronOnCreateAuthorInBatch(t *testing.T) {
	authorDbMock := new(authorrepo.AuthorRepositoryMock)

	authorDbMock.On("CreateAuthorInBatch", mock.Anything, 6).Return(errors.New("error occurred"))

	authorServiceTest := AuthorService{authorDb: authorDbMock}

	resp, err := authorServiceTest.ImportAuthorsFromCSVFile("../../data/authorsreduced.csv")

	require.Error(t, err)
	require.Equal(t, 0, resp)
}

func TestImportAuthorsFromCSVFileError(t *testing.T) {
	authorDbMock := new(authorrepo.AuthorRepositoryMock)

	authorServiceTest := AuthorService{authorDb: authorDbMock}

	_, err := authorServiceTest.ImportAuthorsFromCSVFile("anyfile")
	require.Error(t, err)
}
