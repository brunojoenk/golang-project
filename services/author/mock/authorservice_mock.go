package services

import (
	"github/brunojoenk/golang-test/models/dtos"

	"github.com/stretchr/testify/mock"
)

type AuthorServiceyMock struct {
	mock.Mock
}

func (m *AuthorServiceyMock) GetAllAuthors(filter dtos.GetAuthorsFilter) (*dtos.AuthorResponseMetadata, error) {
	args := m.Called(filter)
	return args.Get(0).(*dtos.AuthorResponseMetadata), args.Error(1)
}

func (m *AuthorServiceyMock) ImportAuthorsFromCSVFile(file string) (int, error) {
	args := m.Called(file)
	return args.Get(0).(int), args.Error(1)
}
