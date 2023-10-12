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

var errGeneric = errors.New("generic error")

func TestGetAllAuthors(t *testing.T) {
	tests := map[string]struct {
		authors                   []entities.Author
		expectedErrorOnGetAuthors error
		expectedErrorResponse     error
	}{
		"success on get all authors": {
			authors: []entities.Author{{Id: 5, Name: "Joenk"}},
		},
		"error occurred on get all authors": {
			authors:                   []entities.Author{{Id: 5, Name: "Joenk"}},
			expectedErrorOnGetAuthors: errGeneric,
			expectedErrorResponse:     errGeneric,
		},
	}
	for testName, tc := range tests {
		t.Run(testName, func(t *testing.T) {
			authorDbMock := new(authorrepomock.AuthorRepositoryMock)

			filter := dtos.GetAuthorsFilter{Pagination: dtos.Pagination{Page: 1, Limit: 10}}
			authorDbMock.On("GetAllAuthors", filter).Return(tc.authors, tc.expectedErrorOnGetAuthors)

			authorServiceTest := authorService{authorDb: authorDbMock}

			resp, err := authorServiceTest.GetAllAuthors(filter)
			if tc.expectedErrorResponse != nil {
				require.Error(t, err)
				require.Equal(t, tc.expectedErrorResponse, err)
			} else {
				require.NoError(t, err)
				require.Nil(t, deep.Equal([]dtos.AuthorResponse{{Id: 5, Name: "Joenk"}}, resp.Authors))
			}
		})
	}
}

func TestImportAuthorsFromCSVFile2(t *testing.T) {
	filePath := "../../data/authorsreduced.csv"
	tests := map[string]struct {
		filePath                         string
		totalAuthorsExpected             int
		expectedErrorCreateAuthorInBatch error
		expectedErrorResponse            error
	}{
		"success on import all authors": {
			filePath:             filePath,
			totalAuthorsExpected: 6,
		},
		"error occurred on import all authors (create author in batch)": {
			filePath:                         filePath,
			totalAuthorsExpected:             0,
			expectedErrorCreateAuthorInBatch: errGeneric,
			expectedErrorResponse:            errGeneric,
		},
	}
	for testName, tc := range tests {
		t.Run(testName, func(t *testing.T) {
			authorDbMock := new(authorrepomock.AuthorRepositoryMock)

			authorDbMock.On("CreateAuthorInBatch", mock.Anything, 6).Return(tc.expectedErrorCreateAuthorInBatch)

			authorServiceTest := authorService{authorDb: authorDbMock}

			resp, err := authorServiceTest.ImportAuthorsFromCSVFile(tc.filePath)
			if tc.expectedErrorResponse != nil {
				require.Equal(t, tc.expectedErrorResponse, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.totalAuthorsExpected, resp)
			}
		})
	}
}

func TestImportAuthorsFromCSVFileError(t *testing.T) {
	authorDbMock := new(authorrepomock.AuthorRepositoryMock)

	authorServiceTest := authorService{authorDb: authorDbMock}

	_, err := authorServiceTest.ImportAuthorsFromCSVFile("anyfile")
	require.Error(t, err)
}
