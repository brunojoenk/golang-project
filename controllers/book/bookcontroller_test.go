package controllers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github/brunojoenk/golang-test/models/dtos"
	"github/brunojoenk/golang-test/utils"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var bookControllerTest *BookController

func init() {
	db, _, _ := sqlmock.New()

	dialector := postgres.New(postgres.Config{
		DSN:                  "sqlmock_db_0",
		DriverName:           "postgres",
		Conn:                 db,
		PreferSimpleProtocol: true,
	})
	gormDB, _ := gorm.Open(dialector, &gorm.Config{})

	bookControllerTest = NewBookController(gormDB)
}

func TestCreateBook(t *testing.T) {
	bookControllerTest.createBookRepo = func(bookRequestCreate dtos.BookRequestCreate) error {
		return nil
	}

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := bookControllerTest.CreateBook(c)

	require.NoError(t, err)
	require.Equal(t, http.StatusCreated, rec.Code)

}

func TestCreateBookErrorOnBody(t *testing.T) {

	jsonStr := `{name: 5}`
	jsonBody, _ := json.Marshal(jsonStr)

	request, err := http.NewRequest("POST", "/books", bytes.NewBuffer(jsonBody))
	recorder := httptest.NewRecorder()
	e := echo.New()
	e.POST("/books", bookControllerTest.CreateBook)
	e.ServeHTTP(recorder, request)

	require.Equal(t, http.StatusBadRequest, recorder.Code)

	if err != nil {
		t.Errorf("Expected nil, received %s", err.Error())
	}

}

func TestCreateBookErrorOnService(t *testing.T) {
	errExpected := errors.New("error occurred")
	bookControllerTest.createBookRepo = func(bookRequestCreate dtos.BookRequestCreate) error {
		return errExpected
	}

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	bookControllerTest.CreateBook(c)

	require.Equal(t, http.StatusInternalServerError, rec.Code)

}

func TestCreateBookWhenAuthorIdIsNotFound(t *testing.T) {

	bookControllerTest.createBookRepo = func(bookRequestCreate dtos.BookRequestCreate) error {
		return utils.ErrAuthorIdNotFound
	}

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	bookControllerTest.CreateBook(c)

	require.Equal(t, http.StatusBadRequest, rec.Code)

}

func TestGetBook(t *testing.T) {
	var (
		bookId          = 12
		bookName        = "harry"
		bookEdition     = "first"
		publicationYear = 2022
		authorName      = "jk rowling"
	)

	var bookIdCalled int

	bookControllerTest.getBookRepo = func(id int) (*dtos.BookResponse, error) {
		bookIdCalled = id
		return &dtos.BookResponse{Name: bookName, Edition: bookEdition, PublicationYear: publicationYear, Authors: authorName}, nil
	}

	request, err := http.NewRequest("GET", "/book/12", nil)
	recorder := httptest.NewRecorder()
	e := echo.New()
	e.GET("/book/:id", bookControllerTest.GetBook)
	e.ServeHTTP(recorder, request)

	books := dtos.BookResponse{
		Name: bookName, Edition: bookEdition, PublicationYear: publicationYear, Authors: authorName,
	}
	respExpected, _ := json.Marshal(books)
	require.Equal(t, fmt.Sprintf("%s%s", respExpected, "\n"), recorder.Body.String())

	require.NoError(t, err)
	require.Equal(t, http.StatusOK, recorder.Code)
	require.Equal(t, bookId, bookIdCalled)

}

func TestGetBookErrorParameterId(t *testing.T) {
	var (
		bookName        = "harry"
		bookEdition     = "first"
		publicationYear = 2022
		authorName      = "jk rowling"
	)

	bookControllerTest.getBookRepo = func(id int) (*dtos.BookResponse, error) {
		return &dtos.BookResponse{Name: bookName, Edition: bookEdition, PublicationYear: publicationYear, Authors: authorName}, nil
	}

	request, err := http.NewRequest("GET", "/book/a", nil)
	recorder := httptest.NewRecorder()
	e := echo.New()
	e.GET("/book/:id", bookControllerTest.GetBook)
	e.ServeHTTP(recorder, request)

	require.NoError(t, err)
	require.Equal(t, http.StatusBadRequest, recorder.Code)

}

func TestGetBookErrorOnQueryDatabase(t *testing.T) {
	var (
		bookId = 12
	)

	var bookIdCalled int
	errExpected := errors.New("error occurred")

	bookControllerTest.getBookRepo = func(id int) (*dtos.BookResponse, error) {
		bookIdCalled = id
		return nil, errExpected
	}

	request, _ := http.NewRequest("GET", "/book/12", nil)
	recorder := httptest.NewRecorder()
	e := echo.New()
	e.GET("/book/:id", bookControllerTest.GetBook)
	e.ServeHTTP(recorder, request)

	require.Equal(t, http.StatusInternalServerError, recorder.Code)
	require.Equal(t, bookId, bookIdCalled)

}

func TestGetBookErrorWhenAuthorIdNotFound(t *testing.T) {
	var (
		bookId = 12
	)

	var bookIdCalled int

	bookControllerTest.getBookRepo = func(id int) (*dtos.BookResponse, error) {
		bookIdCalled = id
		return nil, utils.ErrBookIdNotFound
	}

	request, _ := http.NewRequest("GET", "/book/12", nil)
	recorder := httptest.NewRecorder()
	e := echo.New()
	e.GET("/book/:id", bookControllerTest.GetBook)
	e.ServeHTTP(recorder, request)

	require.Equal(t, http.StatusBadRequest, recorder.Code)
	require.Equal(t, bookId, bookIdCalled)

}

func TestGetAllBook(t *testing.T) {
	var (
		bookName        = "harry"
		bookEdition     = "first"
		publicationYear = 2022
		authorName      = "jk rowling"
	)

	bookControllerTest.getAllBooksRepo = func(filter dtos.GetBooksFilter) (*dtos.BookResponseMetadata, error) {
		return &dtos.BookResponseMetadata{Books: []dtos.BookResponse{{Name: bookName, Edition: bookEdition, PublicationYear: publicationYear, Authors: authorName}},
			Pagination: dtos.Pagination{Page: 1, Limit: 10}}, nil
	}

	request, err := http.NewRequest("GET", "/books", nil)
	recorder := httptest.NewRecorder()
	e := echo.New()
	e.GET("/books", bookControllerTest.GetAllBooks)
	e.ServeHTTP(recorder, request)

	books := &dtos.BookResponseMetadata{Books: []dtos.BookResponse{{Name: bookName, Edition: bookEdition, PublicationYear: publicationYear, Authors: authorName}}, Pagination: dtos.Pagination{Page: 1, Limit: 10}}
	respExpected, _ := json.Marshal(books)
	require.Equal(t, fmt.Sprintf("%s%s", respExpected, "\n"), recorder.Body.String())

	require.NoError(t, err)
	require.Equal(t, http.StatusOK, recorder.Code)

}

func TestGetAllBookErrorOnFilter(t *testing.T) {
	var (
		bookName        = "harry"
		bookEdition     = "first"
		publicationYear = 2022
		authorName      = "jk rowling"
	)

	bookControllerTest.getAllBooksRepo = func(filter dtos.GetBooksFilter) (*dtos.BookResponseMetadata, error) {
		return &dtos.BookResponseMetadata{Books: []dtos.BookResponse{{Name: bookName, Edition: bookEdition, PublicationYear: publicationYear, Authors: authorName}},
			Pagination: dtos.Pagination{Page: 1, Limit: 10}}, nil
	}

	request, _ := http.NewRequest("GET", "/books?publication_year=joenk", nil)
	recorder := httptest.NewRecorder()
	e := echo.New()
	e.GET("/books", bookControllerTest.GetAllBooks)
	e.ServeHTTP(recorder, request)

	require.Equal(t, http.StatusBadRequest, recorder.Code)

}

func TestGetAllBookErrorOnService(t *testing.T) {
	errExpected := errors.New("error occurred")
	bookControllerTest.getAllBooksRepo = func(filter dtos.GetBooksFilter) (*dtos.BookResponseMetadata, error) {
		return nil, errExpected
	}

	request, _ := http.NewRequest("GET", "/books", nil)
	recorder := httptest.NewRecorder()
	e := echo.New()
	e.GET("/books", bookControllerTest.GetAllBooks)
	e.ServeHTTP(recorder, request)

	require.Equal(t, http.StatusInternalServerError, recorder.Code)

}

func TestDeleteBook(t *testing.T) {
	var (
		bookId = 12
	)

	var bookIdCalled int

	bookControllerTest.deleteBookRepo = func(id int) error {
		bookIdCalled = id
		return nil
	}

	request, err := http.NewRequest("DELETE", "/book/12", nil)
	recorder := httptest.NewRecorder()
	e := echo.New()
	e.DELETE("/book/:id", bookControllerTest.DeleteBook)
	e.ServeHTTP(recorder, request)

	require.NoError(t, err)
	require.Equal(t, http.StatusOK, recorder.Code)
	require.Equal(t, bookId, bookIdCalled)

}

func TestDeleteBookErrorOnInvalidParameterId(t *testing.T) {

	bookControllerTest.deleteBookRepo = func(id int) error {
		return nil
	}

	request, err := http.NewRequest("DELETE", "/book/a", nil)
	recorder := httptest.NewRecorder()
	e := echo.New()
	e.DELETE("/book/:id", bookControllerTest.DeleteBook)
	e.ServeHTTP(recorder, request)

	require.NoError(t, err)
	require.Equal(t, http.StatusBadRequest, recorder.Code)

}

func TestDeleteBookErrorOnService(t *testing.T) {
	var (
		bookId = 12
	)

	var bookIdCalled int
	errExpected := errors.New("error occurred")

	bookControllerTest.deleteBookRepo = func(id int) error {
		bookIdCalled = id
		return errExpected
	}

	request, err := http.NewRequest("DELETE", "/book/12", nil)
	recorder := httptest.NewRecorder()
	e := echo.New()
	e.DELETE("/book/:id", bookControllerTest.DeleteBook)
	e.ServeHTTP(recorder, request)

	require.NoError(t, err)
	require.Equal(t, http.StatusInternalServerError, recorder.Code)
	require.Equal(t, bookId, bookIdCalled)

}

func TestUpdateBook(t *testing.T) {
	var (
		bookId = 12
	)

	var bookIdCalled int

	bookControllerTest.updateBookRepo = func(id int, bookRequestUpdate dtos.BookRequestUpdate) error {
		bookIdCalled = id
		return nil
	}

	request, err := http.NewRequest("PUT", "/book/12", nil)
	recorder := httptest.NewRecorder()
	e := echo.New()
	e.PUT("/book/:id", bookControllerTest.UpdateBook)
	e.ServeHTTP(recorder, request)

	require.NoError(t, err)
	require.Equal(t, http.StatusOK, recorder.Code)
	require.Equal(t, bookId, bookIdCalled)

}

func TestUpdateBookErrorOnParametersIdNotFound(t *testing.T) {

	bookControllerTest.updateBookRepo = func(id int, bookRequestUpdate dtos.BookRequestUpdate) error {
		return nil
	}

	request, err := http.NewRequest("PUT", "/book/a", nil)
	recorder := httptest.NewRecorder()
	e := echo.New()
	e.PUT("/book/:id", bookControllerTest.UpdateBook)
	e.ServeHTTP(recorder, request)

	require.NoError(t, err)
	require.Equal(t, http.StatusBadRequest, recorder.Code)

}

func TestUpdateBookErrorOnBody(t *testing.T) {

	jsonStr := `{name: 5}`
	jsonBody, _ := json.Marshal(jsonStr)

	request, err := http.NewRequest("PUT", "/book/12", bytes.NewBuffer(jsonBody))
	recorder := httptest.NewRecorder()
	e := echo.New()
	e.PUT("/book/:id", bookControllerTest.UpdateBook)
	e.ServeHTTP(recorder, request)

	require.NoError(t, err)
	require.Equal(t, http.StatusBadRequest, recorder.Code)

}

func TestUpdateBookErrorOnService(t *testing.T) {

	bookControllerTest.updateBookRepo = func(id int, bookRequestUpdate dtos.BookRequestUpdate) error {
		return errors.New("error occurred")
	}

	request, err := http.NewRequest("PUT", "/book/12", nil)
	recorder := httptest.NewRecorder()
	e := echo.New()
	e.PUT("/book/:id", bookControllerTest.UpdateBook)
	e.ServeHTTP(recorder, request)

	require.NoError(t, err)
	require.Equal(t, http.StatusInternalServerError, recorder.Code)

}

func TestUpdateBookErrorWhenAuthorIdNotFound(t *testing.T) {

	bookControllerTest.updateBookRepo = func(id int, bookRequestUpdate dtos.BookRequestUpdate) error {
		return utils.ErrAuthorIdNotFound
	}

	request, err := http.NewRequest("PUT", "/book/12", nil)
	recorder := httptest.NewRecorder()
	e := echo.New()
	e.PUT("/book/:id", bookControllerTest.UpdateBook)
	e.ServeHTTP(recorder, request)

	require.NoError(t, err)
	require.Equal(t, http.StatusBadRequest, recorder.Code)

}
