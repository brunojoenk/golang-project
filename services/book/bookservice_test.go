package services

import (
	"errors"
	"github/brunojoenk/golang-test/models/dtos"
	"github/brunojoenk/golang-test/models/entities"
	authorrepomock "github/brunojoenk/golang-test/repository/author/mock"
	bookrepomock "github/brunojoenk/golang-test/repository/book/mock"
	"github/brunojoenk/golang-test/utils"
	"testing"

	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

var errGeneric = errors.New("generic error")

func TestCreateBook(t *testing.T) {
	var (
		name            = "book"
		edition         = "edition"
		publicationYear = 2022
		authorId        = 5
		authorName      = "author"
	)
	authors := []entities.Author{{Id: authorId, Name: authorName}}
	book := entities.Book{Name: name, Edition: edition, PublicationYear: publicationYear, Authors: authors}
	tests := map[string]struct {
		book                      entities.Book
		authors                   []entities.Author
		expectedErrorOnGetAuthors error
		expectedErrorOnCreateBook error
		expectedErrorResponse     error
	}{
		"success on create book": {
			book:    book,
			authors: authors,
		},
		"error occurred on create book (author not found)": {
			book:                      book,
			authors:                   authors,
			expectedErrorOnGetAuthors: gorm.ErrRecordNotFound,
			expectedErrorResponse:     utils.ErrAuthorIdNotFound,
		},
		"error occurrente on create book (get authors)": {
			book:                      book,
			authors:                   authors,
			expectedErrorOnGetAuthors: errGeneric,
			expectedErrorResponse:     errGeneric,
		},
		"error occurred on create book (create book)": {
			book:                      book,
			authors:                   authors,
			expectedErrorOnCreateBook: errGeneric,
			expectedErrorResponse:     errGeneric,
		},
	}
	for testName, tc := range tests {
		t.Run(testName, func(t *testing.T) {
			authorDbMock := new(authorrepomock.AuthorRepositoryMock)
			authorDbMock.On("GetAuthor", tc.authors[0].Id).Return(tc.authors[0], tc.expectedErrorOnGetAuthors)

			bookDbMock := new(bookrepomock.BookRepositoryMock)
			bookDbMock.On("CreateBook", tc.book).Return(entities.Book{}, tc.expectedErrorOnCreateBook)

			bookServiceTest := bookService{authorDb: authorDbMock, bookDb: bookDbMock}
			_, err := bookServiceTest.CreateBook(dtos.BookRequestCreate{Name: name, Edition: edition, PublicationYear: publicationYear, Authors: []int{authorId}})
			if tc.expectedErrorResponse != nil {
				require.ErrorIs(t, err, tc.expectedErrorResponse)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestGetAllBooks(t *testing.T) {
	var (
		bookName          = "book"
		edition           = "edition"
		publicationYear   = 2022
		authorId          = 5
		authorName        = "joenk"
		anotherAuthorId   = 7
		anotherAuthorName = "bruno"
		authors           = []entities.Author{{Id: authorId, Name: authorName}, {Id: anotherAuthorId, Name: anotherAuthorName}}
		book1             = entities.Book{Name: bookName, Edition: edition, PublicationYear: publicationYear, Authors: authors}
	)
	tests := map[string]struct {
		booksExpected              []entities.Book
		expectedErrorOnGetAllBooks error
		expectedErrorResponse      error
	}{
		"success on get all books": {
			booksExpected: []entities.Book{book1},
		},
		"error occurred on get all books": {
			expectedErrorOnGetAllBooks: errGeneric,
			expectedErrorResponse:      errGeneric,
		},
	}
	for testName, tc := range tests {
		t.Run(testName, func(t *testing.T) {
			bookDbMock := new(bookrepomock.BookRepositoryMock)
			filter := dtos.GetBooksFilter{Pagination: dtos.Pagination{Page: 1, Limit: 10}}
			bookDbMock.On("GetAllBooks", filter).Return(tc.booksExpected, tc.expectedErrorOnGetAllBooks)

			bookServiceTest := bookService{bookDb: bookDbMock}
			resp, err := bookServiceTest.GetAllBooks(filter)
			if tc.expectedErrorResponse != nil {
				require.ErrorIs(t, err, tc.expectedErrorResponse)
			} else {
				require.NoError(t, err)
				require.Equal(t, resp.Books[0].Name, bookName)
				require.Equal(t, resp.Books[0].Edition, edition)
				require.Equal(t, resp.Books[0].PublicationYear, publicationYear)
				require.Equal(t, resp.Books[0].Authors, "joenk | bruno")
			}
		})
	}

}

func TestDeleteBook(t *testing.T) {
	var (
		bookId = 2
	)
	tests := map[string]struct {
		expectedErrorOnDeleteBook error
		expectedErrorResponse     error
	}{
		"success on delete book": {},
		"error occurred on delete book": {
			expectedErrorOnDeleteBook: errGeneric,
			expectedErrorResponse:     errGeneric,
		},
	}
	for testName, tc := range tests {
		t.Run(testName, func(t *testing.T) {
			bookDbMock := new(bookrepomock.BookRepositoryMock)
			bookDbMock.On("DeleteBook", bookId).Return(tc.expectedErrorOnDeleteBook)

			bookServiceTest := bookService{bookDb: bookDbMock}
			err := bookServiceTest.DeleteBook(bookId)
			if tc.expectedErrorResponse != nil {
				require.ErrorIs(t, err, tc.expectedErrorResponse)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestGetBook(t *testing.T) {
	var (
		bookId            = 1
		bookName          = "book"
		edition           = "edition"
		publicationYear   = 2022
		authorId          = 5
		authorName        = "joenk"
		anotherAuthorId   = 7
		anotherAuthorName = "bruno"
		authors           = []entities.Author{{Id: authorId, Name: authorName}, {Id: anotherAuthorId, Name: anotherAuthorName}}
		book              = entities.Book{Id: bookId, Name: bookName, Edition: edition, PublicationYear: publicationYear, Authors: authors}
	)
	tests := map[string]struct {
		bookExpected           entities.Book
		expectedErrorOnGetBook error
		expectedErrorResponse  error
	}{
		"success on get book": {
			bookExpected: book,
		},
		"error occurred on get book (not found)": {
			expectedErrorOnGetBook: gorm.ErrRecordNotFound,
			expectedErrorResponse:  utils.ErrBookIdNotFound,
		},
		"error occurred on get book": {
			expectedErrorOnGetBook: errGeneric,
			expectedErrorResponse:  errGeneric,
		},
	}
	for testName, tc := range tests {
		t.Run(testName, func(t *testing.T) {
			bookDbMock := new(bookrepomock.BookRepositoryMock)
			bookDbMock.On("GetBook", bookId).Return(tc.bookExpected, tc.expectedErrorOnGetBook)

			bookServiceTest := bookService{bookDb: bookDbMock}
			resp, err := bookServiceTest.GetBook(bookId)
			if tc.expectedErrorResponse != nil {
				require.ErrorIs(t, err, tc.expectedErrorResponse)
			} else {
				require.NoError(t, err)
				require.Equal(t, resp.Name, bookName)
				require.Equal(t, resp.Edition, edition)
				require.Equal(t, resp.PublicationYear, publicationYear)
				require.Equal(t, resp.Authors, "joenk | bruno")
			}
		})
	}
}

func TestUpdateBook(t *testing.T) {
	var (
		bookId          = 5
		bookName        = "book"
		edition         = "edition"
		publicationYear = 2022
		authorId        = 5
		authorName      = "joenk"
		authors         = []entities.Author{{Id: authorId, Name: authorName}}
		book            = entities.Book{Id: bookId, Name: bookName, Edition: edition, PublicationYear: publicationYear, Authors: authors}
	)
	tests := map[string]struct {
		bookExpected             entities.Book
		expectedErrorOnGetBook   error
		expectedErrorOnGetAuthor error
		expectedErrorOnUpdate    error
		expectedErrorResponse    error
	}{
		"success on update book": {
			bookExpected: book,
		},
		"error occurred on update book (get book)": {
			expectedErrorOnGetBook: errGeneric,
			expectedErrorResponse:  errGeneric,
		},
		"error occurred on update book (get author)": {
			expectedErrorOnGetAuthor: errGeneric,
			expectedErrorResponse:    errGeneric,
		},
		"error occurante on update book (author not found)": {
			expectedErrorOnGetAuthor: gorm.ErrRecordNotFound,
			expectedErrorResponse:    utils.ErrAuthorIdNotFound,
		},
		"error occurred on update book (update book)": {
			expectedErrorOnUpdate: errGeneric,
			expectedErrorResponse: errGeneric,
		},
	}
	for testName, tc := range tests {
		t.Run(testName, func(t *testing.T) {
			bookDbMock := new(bookrepomock.BookRepositoryMock)
			bookDbMock.On("GetBook", bookId).Return(book, tc.expectedErrorOnGetBook)
			bookDbMock.On("UpdateBook", book, authors).Return(book, tc.expectedErrorOnUpdate)

			authorDbMock := new(authorrepomock.AuthorRepositoryMock)
			authorDbMock.On("GetAuthor", authorId).Return(authors[0], tc.expectedErrorOnGetAuthor)

			bookServiceTest := bookService{bookDb: bookDbMock, authorDb: authorDbMock}
			_, err := bookServiceTest.UpdateBook(bookId, dtos.BookRequestUpdate{Name: bookName, Edition: edition, PublicationYear: publicationYear, Authors: []int{authorId}})
			if tc.expectedErrorResponse != nil {
				require.ErrorIs(t, err, tc.expectedErrorResponse)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
