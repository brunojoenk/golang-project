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
)

func TestCreateBook(t *testing.T) {
	var (
		name            = "book"
		edition         = "edition"
		publicationYear = 2022
		authorId        = 5
		authorName      = "joenk"
		authors         = []*entities.Author{{Id: authorId, Name: authorName}}
	)

	authorDbMock := new(authorrepomock.AuthorRepositoryMock)
	authorDbMock.On("GetAuthor", authorId).Return(authors[0], nil)

	bookDbMock := new(bookrepomock.BookRepositoryMock)
	bookDbMock.On("CreateBook", &entities.Book{Name: name, Edition: edition, PublicationYear: publicationYear, Authors: authors}).Return(nil)

	bookServiceTest := bookService{authorDb: authorDbMock, bookDb: bookDbMock}
	err := bookServiceTest.CreateBook(dtos.BookRequestCreate{Name: name, Edition: edition, PublicationYear: publicationYear, Authors: []int{authorId}})

	require.NoError(t, err)
}

func TestCreateBook_Error(t *testing.T) {
	var (
		name            = "book"
		edition         = "edition"
		publicationYear = 2022
		authorId        = 5
	)

	errExpected := errors.New("error occurred")

	authorDbMock := new(authorrepomock.AuthorRepositoryMock)
	authorDbMock.On("GetAuthor", authorId).Return(&entities.Author{}, errExpected)

	bookServiceTest := bookService{authorDb: authorDbMock}
	err := bookServiceTest.CreateBook(dtos.BookRequestCreate{Name: name, Edition: edition, PublicationYear: publicationYear, Authors: []int{authorId}})

	require.ErrorIs(t, err, errExpected)
}

func TestCreateBookErrorWhensIsAuthorIdNotFound(t *testing.T) {
	var (
		name            = "book"
		edition         = "edition"
		publicationYear = 2022
		authorId        = 5
	)

	authorDbMock := new(authorrepomock.AuthorRepositoryMock)
	authorDbMock.On("GetAuthor", authorId).Return(&entities.Author{Id: 0}, nil)

	bookServiceTest := bookService{authorDb: authorDbMock}
	err := bookServiceTest.CreateBook(dtos.BookRequestCreate{Name: name, Edition: edition, PublicationYear: publicationYear, Authors: []int{authorId}})

	require.ErrorIs(t, err, utils.ErrAuthorIdNotFound)
}

func TestGetAllBooks(t *testing.T) {
	var (
		bookId            = 1
		bookName          = "book"
		edition           = "edition"
		publicationYear   = 2022
		authorId          = 5
		authorName        = "joenk"
		anotherAuthorId   = 7
		anotherAuthorName = "bruno"
		authors           = []*entities.Author{{Id: authorId, Name: authorName}, {Id: anotherAuthorId, Name: anotherAuthorName}}
	)

	filter := dtos.GetBooksFilter{Pagination: dtos.Pagination{Page: 1, Limit: 10}}
	bookDbMock := new(bookrepomock.BookRepositoryMock)
	bookDbMock.On("GetAllBooks", filter).Return([]entities.Book{{Id: bookId, Name: bookName, Edition: edition, PublicationYear: publicationYear, Authors: authors}}, nil)

	bookServiceTest := bookService{bookDb: bookDbMock}
	resp, err := bookServiceTest.GetAllBooks(filter)

	require.NoError(t, err)
	require.Equal(t, resp.Books[0].Name, bookName)
	require.Equal(t, resp.Books[0].Edition, edition)
	require.Equal(t, resp.Books[0].PublicationYear, publicationYear)
	require.Equal(t, resp.Books[0].Authors, "joenk | bruno")
}

func TestGetAllBooksError(t *testing.T) {
	errExpected := errors.New("error occurred")

	filter := dtos.GetBooksFilter{Pagination: dtos.Pagination{Page: 1, Limit: 10}}
	bookDbMock := new(bookrepomock.BookRepositoryMock)
	bookDbMock.On("GetAllBooks", filter).Return(make([]entities.Book, 0), errExpected)

	bookServiceTest := bookService{bookDb: bookDbMock}
	_, err := bookServiceTest.GetAllBooks(filter)

	require.ErrorIs(t, err, errExpected)
}

func TestDeleteBook(t *testing.T) {
	var (
		bookId = 2
	)

	bookDbMock := new(bookrepomock.BookRepositoryMock)
	bookDbMock.On("DeleteBook", bookId).Return(nil)

	bookServiceTest := bookService{bookDb: bookDbMock}
	err := bookServiceTest.DeleteBook(bookId)

	require.NoError(t, err)
}

func TestDeleteBookError(t *testing.T) {
	var (
		bookId = 2
	)

	errExpected := errors.New("error occurred")

	bookDbMock := new(bookrepomock.BookRepositoryMock)
	bookDbMock.On("DeleteBook", bookId).Return(errExpected)

	bookServiceTest := bookService{bookDb: bookDbMock}
	err := bookServiceTest.DeleteBook(bookId)

	require.Error(t, errExpected, err)
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
		authors           = []*entities.Author{{Id: authorId, Name: authorName}, {Id: anotherAuthorId, Name: anotherAuthorName}}
	)

	bookDbMock := new(bookrepomock.BookRepositoryMock)
	bookDbMock.On("GetBook", bookId).Return(&entities.Book{Id: bookId, Name: bookName, Edition: edition, PublicationYear: publicationYear, Authors: authors}, nil)

	bookServiceTest := bookService{bookDb: bookDbMock}
	resp, err := bookServiceTest.GetBook(bookId)

	require.NoError(t, err)
	require.Equal(t, resp.Name, bookName)
	require.Equal(t, resp.Edition, edition)
	require.Equal(t, resp.PublicationYear, publicationYear)
	require.Equal(t, resp.Authors, "joenk | bruno")
}

func TestGetBookError(t *testing.T) {
	var (
		bookId = 1
	)

	errExpected := errors.New("error occurred")

	bookDbMock := new(bookrepomock.BookRepositoryMock)
	bookDbMock.On("GetBook", bookId).Return(&entities.Book{}, errExpected)

	bookServiceTest := bookService{bookDb: bookDbMock}
	_, err := bookServiceTest.GetBook(bookId)

	require.ErrorIs(t, err, errExpected)
}

func TestGetBookIdNotFound(t *testing.T) {
	var (
		bookId = 1
	)

	bookDbMock := new(bookrepomock.BookRepositoryMock)
	bookDbMock.On("GetBook", bookId).Return(&entities.Book{Id: 0}, nil)

	bookServiceTest := bookService{bookDb: bookDbMock}
	_, err := bookServiceTest.GetBook(bookId)

	require.ErrorIs(t, err, utils.ErrBookIdNotFound)
}

func TestUpdateBook(t *testing.T) {
	var (
		bookId          = 5
		bookName        = "book"
		edition         = "edition"
		publicationYear = 2022
		authorId        = 5
		authorName      = "joenk"
		authors         = []*entities.Author{{Id: authorId, Name: authorName}}
		book            = &entities.Book{Id: bookId, Name: bookName, Edition: edition, PublicationYear: publicationYear, Authors: authors}
	)

	bookDbMock := new(bookrepomock.BookRepositoryMock)
	bookDbMock.On("GetBook", bookId).Return(book, nil)
	bookDbMock.On("UpdateBook", book, authors).Return(nil)

	authorDbMock := new(authorrepomock.AuthorRepositoryMock)
	authorDbMock.On("GetAuthor", authorId).Return(authors[0], nil)

	bookServiceTest := bookService{bookDb: bookDbMock, authorDb: authorDbMock}
	err := bookServiceTest.UpdateBook(bookId, dtos.BookRequestUpdate{Name: bookName, Edition: edition, PublicationYear: publicationYear, Authors: []int{authorId}})

	require.NoError(t, err)
}

func TestUpdateBookErrorOnGetBook(t *testing.T) {
	var (
		bookId          = 5
		bookName        = "book"
		edition         = "edition"
		publicationYear = 2022
		authorId        = 3
	)

	errExpected := errors.New("error occurred")

	bookDbMock := new(bookrepomock.BookRepositoryMock)
	bookDbMock.On("GetBook", bookId).Return(&entities.Book{}, errExpected)

	bookServiceTest := bookService{bookDb: bookDbMock}
	err := bookServiceTest.UpdateBook(bookId, dtos.BookRequestUpdate{Name: bookName, Edition: edition, PublicationYear: publicationYear, Authors: []int{authorId}})

	require.ErrorIs(t, err, errExpected)
}

func TestUpdateBookErrorOnGetAuthor(t *testing.T) {
	var (
		bookId          = 5
		bookName        = "book"
		edition         = "edition"
		publicationYear = 2022
		authorId        = 5
		authorName      = "joenk"
		authors         = []*entities.Author{{Id: authorId, Name: authorName}}
	)

	errExpected := errors.New("error occurred")

	bookDbMock := new(bookrepomock.BookRepositoryMock)
	bookDbMock.On("GetBook", bookId).Return(&entities.Book{Id: bookId, Name: bookName, Edition: edition, PublicationYear: publicationYear, Authors: authors}, nil)

	authorDbMock := new(authorrepomock.AuthorRepositoryMock)
	authorDbMock.On("GetAuthor", authorId).Return(&entities.Author{}, errExpected)

	bookServiceTest := bookService{bookDb: bookDbMock, authorDb: authorDbMock}
	err := bookServiceTest.UpdateBook(bookId, dtos.BookRequestUpdate{Name: bookName, Edition: edition, PublicationYear: publicationYear, Authors: []int{authorId}})

	require.ErrorIs(t, err, errExpected)
}

func TestUpdateBookErrorOnAuthorIdNotFound(t *testing.T) {
	var (
		bookId          = 5
		bookName        = "book"
		edition         = "edition"
		publicationYear = 2022
		authorId        = 5
		authorName      = "joenk"
		authors         = []*entities.Author{{Id: authorId, Name: authorName}}
	)

	bookDbMock := new(bookrepomock.BookRepositoryMock)
	bookDbMock.On("GetBook", bookId).Return(&entities.Book{Id: bookId, Name: bookName, Edition: edition, PublicationYear: publicationYear, Authors: authors}, nil)

	authorDbMock := new(authorrepomock.AuthorRepositoryMock)
	authorDbMock.On("GetAuthor", authorId).Return(&entities.Author{Id: 0}, nil)

	bookServiceTest := bookService{bookDb: bookDbMock, authorDb: authorDbMock}
	err := bookServiceTest.UpdateBook(bookId, dtos.BookRequestUpdate{Name: bookName, Edition: edition, PublicationYear: publicationYear, Authors: []int{authorId}})

	require.ErrorIs(t, err, utils.ErrAuthorIdNotFound)
}

func TestUpdateBookErrorOnUpdate(t *testing.T) {
	var (
		bookId          = 5
		bookName        = "book"
		edition         = "edition"
		publicationYear = 2022
		authorId        = 5
		authorName      = "joenk"
		authors         = []*entities.Author{{Id: authorId, Name: authorName}}
		book            = &entities.Book{Id: bookId, Name: bookName, Edition: edition, PublicationYear: publicationYear, Authors: authors}
	)

	errExpected := errors.New("error occurred")

	bookDbMock := new(bookrepomock.BookRepositoryMock)
	bookDbMock.On("GetBook", bookId).Return(book, nil)
	bookDbMock.On("UpdateBook", book, authors).Return(errExpected)

	authorDbMock := new(authorrepomock.AuthorRepositoryMock)
	authorDbMock.On("GetAuthor", authorId).Return(authors[0], nil)

	bookServiceTest := bookService{bookDb: bookDbMock, authorDb: authorDbMock}
	err := bookServiceTest.UpdateBook(bookId, dtos.BookRequestUpdate{Name: bookName, Edition: edition, PublicationYear: publicationYear, Authors: []int{authorId}})

	require.ErrorIs(t, err, errExpected)
}
