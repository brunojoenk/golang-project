package services

import (
	"errors"
	"github/brunojoenk/golang-test/models/dtos"
	"github/brunojoenk/golang-test/models/entities"
	"testing"

	"github.com/go-test/deep"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	authorrepomock "github/brunojoenk/golang-test/repository/author/mock"
)

func TestGetAllAuthors(t *testing.T) {

	authorDbMock := new(authorrepomock.AuthorRepositoryMock)

	filter := dtos.GetAuthorsFilter{Pagination: dtos.Pagination{Page: 1, Limit: 10}}
	authorDbMock.On("GetAllAuthors", filter).Return([]entities.Author{{Id: 5, Name: "Joenk"}}, nil)

	authorServiceTest := authorService{authorDb: authorDbMock}

	resp, err := authorServiceTest.GetAllAuthors(filter)
	require.NoError(t, err)
	require.Nil(t, deep.Equal([]dtos.AuthorResponse{{Id: 5, Name: "Joenk"}}, resp.Authors))
	require.Equal(t, resp.Pagination.Limit, 10)
	require.Equal(t, resp.Pagination.Page, 1)
}

func TestGetAllAuthorsError(t *testing.T) {

	authorDbMock := new(authorrepomock.AuthorRepositoryMock)

	filter := dtos.GetAuthorsFilter{Pagination: dtos.Pagination{Page: 1, Limit: 10}}
	authorDbMock.On("GetAllAuthors", filter).Return([]entities.Author{}, errors.New("Error on test"))

	authorServiceTest := authorService{authorDb: authorDbMock}

	_, err := authorServiceTest.GetAllAuthors(filter)
	require.Error(t, err)
}

func TestImportAuthorsFromCSVFile(t *testing.T) {
	authorDbMock := new(authorrepomock.AuthorRepositoryMock)

	authorDbMock.On("CreateAuthorInBatch", mock.Anything, 6).Return(nil)

	authorServiceTest := authorService{authorDb: authorDbMock}

	resp, err := authorServiceTest.ImportAuthorsFromCSVFile("../../data/authorsreduced.csv")
	require.NoError(t, err)
	require.Equal(t, 6, resp)
}

func TestImportAuthorsFromCSVFileErronOnCreateAuthorInBatch(t *testing.T) {
	authorDbMock := new(authorrepomock.AuthorRepositoryMock)

	authorDbMock.On("CreateAuthorInBatch", mock.Anything, 6).Return(errors.New("error occurred"))

	authorServiceTest := authorService{authorDb: authorDbMock}

	resp, err := authorServiceTest.ImportAuthorsFromCSVFile("../../data/authorsreduced.csv")

	require.Error(t, err)
	require.Equal(t, 0, resp)
}

func TestImportAuthorsFromCSVFileError(t *testing.T) {
	authorDbMock := new(authorrepomock.AuthorRepositoryMock)

	authorServiceTest := authorService{authorDb: authorDbMock}

	_, err := authorServiceTest.ImportAuthorsFromCSVFile("anyfile")
	require.Error(t, err)
}
