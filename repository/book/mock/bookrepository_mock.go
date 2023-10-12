package repository

import (
	"github/brunojoenk/golang-test/models/dtos"
	"github/brunojoenk/golang-test/models/entities"

	"github.com/stretchr/testify/mock"
)

type BookRepositoryMock struct {
	mock.Mock
}

func (m *BookRepositoryMock) CreateBook(book entities.Book) error {
	args := m.Called(book)
	return args.Error(0)
}

func (m *BookRepositoryMock) UpdateBook(book entities.Book, authors []entities.Author) error {
	args := m.Called(book, authors)
	return args.Error(0)
}

func (m *BookRepositoryMock) GetBook(id int) (entities.Book, error) {
	args := m.Called(id)
	return args.Get(0).(entities.Book), args.Error(1)
}

func (m *BookRepositoryMock) GetAllBooks(filter dtos.GetBooksFilter) ([]entities.Book, error) {
	args := m.Called(filter)
	return args.Get(0).([]entities.Book), args.Error(1)
}

func (m *BookRepositoryMock) DeleteBook(id int) error {
	args := m.Called(id)
	return args.Error(0)
}
