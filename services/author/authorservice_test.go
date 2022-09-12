package services

import (
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
