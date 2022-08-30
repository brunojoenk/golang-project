package services

import (
	"errors"
	"github/brunojoenk/golang-test/models/dtos"
	"github/brunojoenk/golang-test/models/entities"
	"github/brunojoenk/golang-test/utils"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var bookServiceTest *BookService

func init() {
	db, _, _ := sqlmock.New()

	dialector := postgres.New(postgres.Config{
		DSN:                  "sqlmock_db_0",
		DriverName:           "postgres",
		Conn:                 db,
		PreferSimpleProtocol: true,
	})
	gormDB, _ := gorm.Open(dialector, &gorm.Config{})

	bookServiceTest = NewBookService(gormDB)
}

func TestCreateBook(t *testing.T) {
	var (
		name            = "book"
		edition         = "edition"
		publicationYear = 2022
		authorId        = 5
		authorName      = "joenk"
		authors         = []*entities.Author{{Id: authorId, Name: authorName}}
	)

	bookServiceTest.getAuthorRepo = func(id int) (*entities.Author, error) {
		return authors[0], nil
	}

	var bookToCreate *entities.Book

	bookServiceTest.createBookRepo = func(book *entities.Book) error {
		bookToCreate = book
		return nil
	}

	err := bookServiceTest.CreateBook(dtos.BookRequestCreateUpdate{Name: name, Edition: edition, PublicationYear: publicationYear, Authors: []int{authorId}})

	require.NoError(t, err)
	require.Equal(t, bookToCreate.Name, name)
	require.Equal(t, bookToCreate.Edition, edition)
	require.Equal(t, bookToCreate.PublicationYear, publicationYear)
	require.Equal(t, bookToCreate.Authors, authors)
}

func TestCreateBook_Error(t *testing.T) {
	var (
		name            = "book"
		edition         = "edition"
		publicationYear = 2022
		authorId        = 5
	)

	errExpected := errors.New("error occurred")

	bookServiceTest.getAuthorRepo = func(id int) (*entities.Author, error) {
		return nil, errExpected
	}

	err := bookServiceTest.CreateBook(dtos.BookRequestCreateUpdate{Name: name, Edition: edition, PublicationYear: publicationYear, Authors: []int{authorId}})

	require.ErrorIs(t, err, errExpected)
}

func TestCreateBook_Error_Whens_Is_Author_Id_Not_Found(t *testing.T) {
	var (
		name            = "book"
		edition         = "edition"
		publicationYear = 2022
		authorId        = 5
	)

	bookServiceTest.getAuthorRepo = func(id int) (*entities.Author, error) {
		return &entities.Author{Id: 0}, nil
	}

	err := bookServiceTest.CreateBook(dtos.BookRequestCreateUpdate{Name: name, Edition: edition, PublicationYear: publicationYear, Authors: []int{authorId}})

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

	bookServiceTest.getAllBooksRepo = func(filter dtos.GetBooksFilter) ([]entities.Book, error) {
		return []entities.Book{{Id: bookId, Name: bookName, Edition: edition, PublicationYear: publicationYear, Authors: authors}}, nil
	}

	resp, err := bookServiceTest.GetAllBooks(dtos.GetBooksFilter{})

	require.NoError(t, err)
	require.Equal(t, resp.Books[0].Name, bookName)
	require.Equal(t, resp.Books[0].Edition, edition)
	require.Equal(t, resp.Books[0].PublicationYear, publicationYear)
	require.Equal(t, resp.Books[0].Authors, "joenk | bruno")
}

func TestGetAllBooksError(t *testing.T) {
	errExpected := errors.New("error occurred")

	bookServiceTest.getAllBooksRepo = func(filter dtos.GetBooksFilter) ([]entities.Book, error) {
		return nil, errExpected
	}

	_, err := bookServiceTest.GetAllBooks(dtos.GetBooksFilter{})

	require.ErrorIs(t, err, errExpected)
}

func TestDeleteBook(t *testing.T) {
	var (
		bookId = 2
	)

	var bookIdCalledExpected int

	bookServiceTest.deleteBookRepo = func(id int) error {
		bookIdCalledExpected = id
		return nil
	}

	err := bookServiceTest.DeleteBook(bookId)

	require.NoError(t, err)
	require.Equal(t, bookId, bookIdCalledExpected)
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

	var bookIdCalledExpected int

	bookServiceTest.getBookRepo = func(id int) (*entities.Book, error) {
		bookIdCalledExpected = id
		return &entities.Book{Id: bookId, Name: bookName, Edition: edition, PublicationYear: publicationYear, Authors: authors}, nil
	}

	resp, err := bookServiceTest.GetBook(bookId)

	require.NoError(t, err)
	require.Equal(t, bookId, bookIdCalledExpected)
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

	bookServiceTest.getBookRepo = func(id int) (*entities.Book, error) {
		return nil, errExpected
	}

	_, err := bookServiceTest.GetBook(bookId)

	require.ErrorIs(t, err, errExpected)
}

func TestGetBookIdNotFound(t *testing.T) {
	var (
		bookId = 1
	)

	bookServiceTest.getBookRepo = func(id int) (*entities.Book, error) {
		return &entities.Book{Id: 0}, nil
	}

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
	)

	var bookIdCalledExpected int

	bookServiceTest.getBookRepo = func(id int) (*entities.Book, error) {
		bookIdCalledExpected = id
		return &entities.Book{Id: bookId, Name: bookName, Edition: edition, PublicationYear: publicationYear, Authors: authors}, nil
	}

	bookServiceTest.getAuthorRepo = func(id int) (*entities.Author, error) {
		return authors[0], nil
	}

	var bookToUpdate *entities.Book

	bookServiceTest.updateBookRepo = func(book *entities.Book, authors []*entities.Author) error {
		bookToUpdate = book
		return nil
	}

	err := bookServiceTest.UpdateBook(bookId, dtos.BookRequestCreateUpdate{Name: bookName, Edition: edition, PublicationYear: publicationYear, Authors: []int{authorId}})

	require.NoError(t, err)
	require.Equal(t, bookIdCalledExpected, bookId)
	require.Equal(t, bookToUpdate.Name, bookName)
	require.Equal(t, bookToUpdate.Edition, edition)
	require.Equal(t, bookToUpdate.PublicationYear, publicationYear)
	require.Equal(t, bookToUpdate.Authors, authors)
}

func TestUpdateBookErrorOnGetBook(t *testing.T) {
	var (
		bookId          = 5
		bookName        = "book"
		edition         = "edition"
		publicationYear = 2022
		authorId        = 3
	)

	var bookIdCalledExpected int

	errExpected := errors.New("error occurred")

	bookServiceTest.getBookRepo = func(id int) (*entities.Book, error) {
		bookIdCalledExpected = id
		return nil, errExpected
	}

	err := bookServiceTest.UpdateBook(bookId, dtos.BookRequestCreateUpdate{Name: bookName, Edition: edition, PublicationYear: publicationYear, Authors: []int{authorId}})

	require.ErrorIs(t, err, errExpected)
	require.Equal(t, bookId, bookIdCalledExpected)
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

	var bookIdCalledExpected int

	bookServiceTest.getBookRepo = func(id int) (*entities.Book, error) {
		bookIdCalledExpected = id
		return &entities.Book{Id: bookId, Name: bookName, Edition: edition, PublicationYear: publicationYear, Authors: authors}, nil
	}

	errExpected := errors.New("error occurred")

	bookServiceTest.getAuthorRepo = func(id int) (*entities.Author, error) {
		return nil, errExpected
	}

	err := bookServiceTest.UpdateBook(bookId, dtos.BookRequestCreateUpdate{Name: bookName, Edition: edition, PublicationYear: publicationYear, Authors: []int{authorId}})

	require.ErrorIs(t, err, errExpected)
	require.Equal(t, bookIdCalledExpected, bookId)
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

	var bookIdCalledExpected int

	bookServiceTest.getBookRepo = func(id int) (*entities.Book, error) {
		bookIdCalledExpected = id
		return &entities.Book{Id: bookId, Name: bookName, Edition: edition, PublicationYear: publicationYear, Authors: authors}, nil
	}

	bookServiceTest.getAuthorRepo = func(id int) (*entities.Author, error) {
		return &entities.Author{Id: 0}, nil
	}

	err := bookServiceTest.UpdateBook(bookId, dtos.BookRequestCreateUpdate{Name: bookName, Edition: edition, PublicationYear: publicationYear, Authors: []int{authorId}})

	require.ErrorIs(t, err, utils.ErrAuthorIdNotFound)
	require.Equal(t, bookIdCalledExpected, bookId)
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
	)

	var bookIdCalledExpected int

	bookServiceTest.getBookRepo = func(id int) (*entities.Book, error) {
		bookIdCalledExpected = id
		return &entities.Book{Id: bookId, Name: bookName, Edition: edition, PublicationYear: publicationYear, Authors: authors}, nil
	}

	bookServiceTest.getAuthorRepo = func(id int) (*entities.Author, error) {
		return authors[0], nil
	}

	errExpected := errors.New("error occurred")

	bookServiceTest.updateBookRepo = func(book *entities.Book, authors []*entities.Author) error {
		return errExpected
	}

	err := bookServiceTest.UpdateBook(bookId, dtos.BookRequestCreateUpdate{Name: bookName, Edition: edition, PublicationYear: publicationYear, Authors: []int{authorId}})

	require.ErrorIs(t, err, errExpected)
	require.Equal(t, bookIdCalledExpected, bookId)
}
