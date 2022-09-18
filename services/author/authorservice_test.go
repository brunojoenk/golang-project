package services

import (
	"errors"
	"github/brunojoenk/golang-test/models/dtos"
	"github/brunojoenk/golang-test/models/entities"
	"testing"

	"github.com/go-test/deep"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type authorDbMock struct {
	mock.Mock
}

func (m *authorDbMock) GetAllAuthors(filter dtos.GetAuthorsFilter) ([]entities.Author, error) {
	args := m.Called(filter)
	return args.Get(0).([]entities.Author), args.Error(1)
}

func (m *authorDbMock) CreateAuthorInBatch(author []entities.Author, batchSize int) error {
	args := m.Called(author, batchSize)
	return args.Error(0)
}

func (m *authorDbMock) GetAuthor(id int) (*entities.Author, error) {
	args := m.Called(id)
	return args.Get(0).(*entities.Author), args.Error(1)
}

func TestGetAllAuthors(t *testing.T) {

	authorDbMock := new(authorDbMock)

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

	authorDbMock := new(authorDbMock)

	filter := dtos.GetAuthorsFilter{Pagination: dtos.Pagination{Page: 1, Limit: 10}}
	authorDbMock.On("GetAllAuthors", filter).Return([]entities.Author{}, errors.New("Error on test"))

	authorServiceTest := AuthorService{authorDb: authorDbMock}

	_, err := authorServiceTest.GetAllAuthors(filter)
	require.Error(t, err)
}

func TestImportAuthorsFromCSVFile(t *testing.T) {
	authorDbMock := new(authorDbMock)

	authorDbMock.On("CreateAuthorInBatch", mock.Anything, 6).Return(nil)

	authorServiceTest := AuthorService{authorDb: authorDbMock}

	resp, err := authorServiceTest.ImportAuthorsFromCSVFile("../../data/authorsreduced.csv")
	require.NoError(t, err)
	require.Equal(t, 6, resp)
}

func TestImportAuthorsFromCSVFileErronOnCreateAuthorInBatch(t *testing.T) {
	authorDbMock := new(authorDbMock)

	authorDbMock.On("CreateAuthorInBatch", mock.Anything, 6).Return(errors.New("error occurred"))

	authorServiceTest := AuthorService{authorDb: authorDbMock}

	resp, err := authorServiceTest.ImportAuthorsFromCSVFile("../../data/authorsreduced.csv")

	require.Error(t, err)
	require.Equal(t, 0, resp)
}

func TestImportAuthorsFromCSVFileError(t *testing.T) {
	authorDbMock := new(authorDbMock)

	authorServiceTest := AuthorService{authorDb: authorDbMock}

	_, err := authorServiceTest.ImportAuthorsFromCSVFile("anyfile")
	require.Error(t, err)
}
