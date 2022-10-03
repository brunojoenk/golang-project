package repository

import (
	"github/brunojoenk/golang-test/models/dtos"
	"github/brunojoenk/golang-test/models/entities"

	"github.com/stretchr/testify/mock"
)

type AuthorRepositoryMock struct {
	mock.Mock
}

func (m *AuthorRepositoryMock) GetAllAuthors(filter dtos.GetAuthorsFilter) ([]entities.Author, error) {
	args := m.Called(filter)
	return args.Get(0).([]entities.Author), args.Error(1)
}

func (m *AuthorRepositoryMock) CreateAuthorInBatch(author []entities.Author, batchSize int) error {
	args := m.Called(author, batchSize)
	return args.Error(0)
}

func (m *AuthorRepositoryMock) GetAuthor(id int) (*entities.Author, error) {
	args := m.Called(id)
	return args.Get(0).(*entities.Author), args.Error(1)
}
