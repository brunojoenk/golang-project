package services

import (
	"github/brunojoenk/golang-test/models/dtos"

	"github.com/stretchr/testify/mock"
)

/*
	CreateBook(bookRequestCreate dtos.BookRequestCreate) error
	GetAllBooks(filter dtos.GetBooksFilter) (*dtos.BookResponseMetadata, error)
	DeleteBook(id int) error
	GetBook(id int) (*dtos.BookResponse, error)
	UpdateBook(id int, bookRequestUpdate dtos.BookRequestUpdate) error
*/

type BookServiceMock struct {
	mock.Mock
}

func (m *BookServiceMock) CreateBook(bookRequestCreate dtos.BookRequestCreate) error {
	args := m.Called(bookRequestCreate)
	return args.Error(0)
}

func (m *BookServiceMock) GetAllBooks(filter dtos.GetBooksFilter) (dtos.BookResponseMetadata, error) {
	args := m.Called(filter)
	return args.Get(0).(dtos.BookResponseMetadata), args.Error(1)
}

func (m *BookServiceMock) DeleteBook(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *BookServiceMock) GetBook(id int) (dtos.BookResponse, error) {
	args := m.Called(id)
	return args.Get(0).(dtos.BookResponse), args.Error(1)
}

func (m *BookServiceMock) UpdateBook(id int, bookRequestUpdate dtos.BookRequestUpdate) error {
	args := m.Called(id, bookRequestUpdate)
	return args.Error(0)
}
