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
	"strings"
	"testing"

	bookservicemock "github/brunojoenk/golang-test/services/book/mock"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
)

func TestCreateBook(t *testing.T) {
	bookServiceMock := new(bookservicemock.BookServiceMock)
	bookServiceMock.On("CreateBook", dtos.BookRequestCreate{}).Return(dtos.BookResponse{Id: 1}, nil)

	bookControllerTest := bookController{bookService: bookServiceMock}

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/book", nil)
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
	bookControllerTest := bookController{}
	e.POST("/books", bookControllerTest.CreateBook)
	e.ServeHTTP(recorder, request)

	require.Equal(t, http.StatusBadRequest, recorder.Code)

	if err != nil {
		t.Errorf("Expected nil, received %s", err.Error())
	}

}

func TestCreateBookErrorOnService(t *testing.T) {
	errExpected := errors.New("error occurred")
	bookServiceMock := new(bookservicemock.BookServiceMock)
	bookServiceMock.On("CreateBook", dtos.BookRequestCreate{}).Return(dtos.BookResponse{}, errExpected)

	bookControllerTest := bookController{bookService: bookServiceMock}

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	bookControllerTest.CreateBook(c)

	require.Equal(t, http.StatusInternalServerError, rec.Code)

}

func TestCreateBookWhenAuthorIdIsNotFound(t *testing.T) {
	bookServiceMock := new(bookservicemock.BookServiceMock)
	bookServiceMock.On("CreateBook", dtos.BookRequestCreate{}).Return(dtos.BookResponse{}, utils.ErrAuthorIdNotFound)

	bookControllerTest := bookController{bookService: bookServiceMock}

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

	bookServiceMock := new(bookservicemock.BookServiceMock)
	bookResponse := dtos.BookResponse{Id: bookId, Name: bookName, Edition: bookEdition, PublicationYear: publicationYear, Authors: authorName}
	bookServiceMock.On("GetBook", bookId).Return(bookResponse, nil)

	bookControllerTest := bookController{bookService: bookServiceMock}

	request, err := http.NewRequest("GET", "/book/12", nil)
	recorder := httptest.NewRecorder()
	e := echo.New()
	e.GET("/book/:id", bookControllerTest.GetBook)
	e.ServeHTTP(recorder, request)

	books := dtos.BookResponse{
		Id: bookId, Name: bookName, Edition: bookEdition, PublicationYear: publicationYear, Authors: authorName,
	}
	respExpected, _ := json.Marshal(books)
	require.Equal(t, fmt.Sprintf("%s%s", respExpected, "\n"), recorder.Body.String())

	require.NoError(t, err)
	require.Equal(t, http.StatusOK, recorder.Code)

}

func TestGetBookErrorParameterId(t *testing.T) {
	bookServiceMock := new(bookservicemock.BookServiceMock)
	bookControllerTest := bookController{bookService: bookServiceMock}

	request, err := http.NewRequest("GET", "/book/a", nil)
	recorder := httptest.NewRecorder()
	e := echo.New()
	e.GET("/book/:id", bookControllerTest.GetBook)
	e.ServeHTTP(recorder, request)

	require.NoError(t, err)
	require.Equal(t, http.StatusBadRequest, recorder.Code)

}

func TestGetBookErrorOnQueryDatabase(t *testing.T) {
	errExpected := errors.New("error occurred")
	bookId := 12

	bookServiceMock := new(bookservicemock.BookServiceMock)
	bookServiceMock.On("GetBook", bookId).Return(dtos.BookResponse{}, errExpected)

	bookControllerTest := bookController{bookService: bookServiceMock}

	request, _ := http.NewRequest("GET", fmt.Sprintf("/book/%v", bookId), nil)
	recorder := httptest.NewRecorder()
	e := echo.New()
	e.GET("/book/:id", bookControllerTest.GetBook)
	e.ServeHTTP(recorder, request)

	require.Equal(t, http.StatusInternalServerError, recorder.Code)

}

func TestGetBookErrorWhenBookIdNotFound(t *testing.T) {
	bookId := 12

	bookServiceMock := new(bookservicemock.BookServiceMock)
	bookServiceMock.On("GetBook", bookId).Return(dtos.BookResponse{}, utils.ErrBookIdNotFound)

	bookControllerTest := bookController{bookService: bookServiceMock}

	request, _ := http.NewRequest("GET", fmt.Sprintf("/book/%v", bookId), nil)
	recorder := httptest.NewRecorder()
	e := echo.New()
	e.GET("/book/:id", bookControllerTest.GetBook)
	e.ServeHTTP(recorder, request)

	require.Equal(t, http.StatusNotFound, recorder.Code)

}

func TestGetAllBook(t *testing.T) {
	var (
		bookId          = 12
		bookName        = "harry"
		bookEdition     = "first"
		publicationYear = 2022
		authorName      = "jk rowling"
	)

	booksResponse := dtos.BookResponseMetadata{Books: []dtos.BookResponse{{Id: bookId, Name: bookName, Edition: bookEdition, PublicationYear: publicationYear, Authors: authorName}},
		Pagination: dtos.Pagination{Page: 1, Limit: 10}}

	bookServiceMock := new(bookservicemock.BookServiceMock)
	bookServiceMock.On("GetAllBooks", dtos.GetBooksFilter{}).Return(booksResponse, nil)

	bookControllerTest := bookController{bookService: bookServiceMock}

	request, err := http.NewRequest("GET", "/books", nil)
	recorder := httptest.NewRecorder()
	e := echo.New()
	e.GET("/books", bookControllerTest.GetAllBooks)
	e.ServeHTTP(recorder, request)

	respExpected, _ := json.Marshal(booksResponse)
	require.Equal(t, fmt.Sprintf("%s%s", respExpected, "\n"), recorder.Body.String())

	require.NoError(t, err)
	require.Equal(t, http.StatusOK, recorder.Code)

}

func TestGetAllBookErrorOnFilter(t *testing.T) {
	bookControllerTest := bookController{}

	request, _ := http.NewRequest("GET", "/books?publication_year=joenk", nil)
	recorder := httptest.NewRecorder()
	e := echo.New()
	e.GET("/books", bookControllerTest.GetAllBooks)
	e.ServeHTTP(recorder, request)

	require.Equal(t, http.StatusBadRequest, recorder.Code)

}

func TestGetAllBookErrorOnService(t *testing.T) {
	errExpected := errors.New("error occurred")

	bookServiceMock := new(bookservicemock.BookServiceMock)
	bookServiceMock.On("GetAllBooks", dtos.GetBooksFilter{}).Return(dtos.BookResponseMetadata{}, errExpected)

	bookControllerTest := bookController{bookService: bookServiceMock}

	request, _ := http.NewRequest("GET", "/books", nil)
	recorder := httptest.NewRecorder()
	e := echo.New()
	e.GET("/books", bookControllerTest.GetAllBooks)
	e.ServeHTTP(recorder, request)

	require.Equal(t, http.StatusInternalServerError, recorder.Code)

}

func TestDeleteBook(t *testing.T) {
	bookId := 12

	bookServiceMock := new(bookservicemock.BookServiceMock)
	bookServiceMock.On("DeleteBook", bookId).Return(nil)

	bookControllerTest := bookController{bookService: bookServiceMock}

	request, err := http.NewRequest("DELETE", fmt.Sprintf("/book/%v", bookId), nil)
	recorder := httptest.NewRecorder()
	e := echo.New()
	e.DELETE("/book/:id", bookControllerTest.DeleteBook)
	e.ServeHTTP(recorder, request)

	require.NoError(t, err)
	require.Equal(t, http.StatusNoContent, recorder.Code)

}

func TestDeleteBookErrorOnInvalidParameterId(t *testing.T) {
	bookControllerTest := bookController{}

	request, err := http.NewRequest("DELETE", "/book/a", nil)
	recorder := httptest.NewRecorder()
	e := echo.New()
	e.DELETE("/book/:id", bookControllerTest.DeleteBook)
	e.ServeHTTP(recorder, request)

	require.NoError(t, err)
	require.Equal(t, http.StatusBadRequest, recorder.Code)

}

func TestDeleteBookErrorOnService(t *testing.T) {
	bookId := 12

	errExpected := errors.New("error occurred")

	bookServiceMock := new(bookservicemock.BookServiceMock)
	bookServiceMock.On("DeleteBook", bookId).Return(errExpected)

	bookControllerTest := bookController{bookService: bookServiceMock}

	request, err := http.NewRequest("DELETE", fmt.Sprintf("/book/%v", bookId), nil)
	recorder := httptest.NewRecorder()
	e := echo.New()
	e.DELETE("/book/:id", bookControllerTest.DeleteBook)
	e.ServeHTTP(recorder, request)

	require.NoError(t, err)
	require.Equal(t, http.StatusInternalServerError, recorder.Code)

}

func TestUpdateBook(t *testing.T) {
	bookId := 12

	bodyRequest := strings.NewReader(`{"name":"Harry Potter 2","edition":"Segunda edição","publication_year":2022,"authors":[5]}`)
	bookRequestUpdate := dtos.BookRequestUpdate{Name: "Harry Potter 2", Edition: "Segunda edição", PublicationYear: 2022, Authors: []int{5}}

	bookServiceMock := new(bookservicemock.BookServiceMock)
	bookServiceMock.On("UpdateBook", bookId, bookRequestUpdate).Return(dtos.BookResponse{}, nil)

	bookControllerTest := bookController{bookService: bookServiceMock}

	request, err := http.NewRequest("PUT", fmt.Sprintf("/book/%v", bookId), bodyRequest)
	request.Header.Add("Content-type", "application/json")

	recorder := httptest.NewRecorder()
	e := echo.New()
	e.PUT("/book/:id", bookControllerTest.UpdateBook)
	e.ServeHTTP(recorder, request)

	require.NoError(t, err)
	require.Equal(t, http.StatusOK, recorder.Code)

}

func TestUpdateBookErrorOnParametersIdNotFound(t *testing.T) {
	bookControllerTest := bookController{}

	request, err := http.NewRequest("PUT", "/book/a", nil)
	recorder := httptest.NewRecorder()
	e := echo.New()
	e.PUT("/book/:id", bookControllerTest.UpdateBook)
	e.ServeHTTP(recorder, request)

	require.NoError(t, err)
	require.Equal(t, http.StatusBadRequest, recorder.Code)

}

func TestUpdateBookErrorOnBody(t *testing.T) {
	bookControllerTest := bookController{}

	bodyRequest := strings.NewReader(`{name: 5}`)

	request, err := http.NewRequest("PUT", "/book/12", bodyRequest)
	request.Header.Add("Content-type", "application/json")

	recorder := httptest.NewRecorder()
	e := echo.New()
	e.PUT("/book/:id", bookControllerTest.UpdateBook)
	e.ServeHTTP(recorder, request)

	require.NoError(t, err)
	require.Equal(t, http.StatusBadRequest, recorder.Code)

}

func TestUpdateBookErrorOnService(t *testing.T) {
	bookId := 12

	bodyRequest := strings.NewReader(`{"name":"Harry Potter 2","edition":"Segunda edição","publication_year":2022,"authors":[5]}`)
	bookRequestUpdate := dtos.BookRequestUpdate{Name: "Harry Potter 2", Edition: "Segunda edição", PublicationYear: 2022, Authors: []int{5}}

	bookServiceMock := new(bookservicemock.BookServiceMock)
	bookServiceMock.On("UpdateBook", bookId, bookRequestUpdate).Return(dtos.BookResponse{}, errors.New("error occurred"))

	bookControllerTest := bookController{bookService: bookServiceMock}

	request, err := http.NewRequest("PUT", fmt.Sprintf("/book/%v", bookId), bodyRequest)
	request.Header.Add("Content-type", "application/json")

	recorder := httptest.NewRecorder()
	e := echo.New()
	e.PUT("/book/:id", bookControllerTest.UpdateBook)
	e.ServeHTTP(recorder, request)

	require.NoError(t, err)
	require.Equal(t, http.StatusInternalServerError, recorder.Code)

}

func TestUpdateBookErrorWhenAuthorIdNotFound(t *testing.T) {
	bookId := 12

	bodyRequest := strings.NewReader(`{"name":"Harry Potter 2","edition":"Segunda edição","publication_year":2022,"authors":[5]}`)
	bookRequestUpdate := dtos.BookRequestUpdate{Name: "Harry Potter 2", Edition: "Segunda edição", PublicationYear: 2022, Authors: []int{5}}

	bookServiceMock := new(bookservicemock.BookServiceMock)
	bookServiceMock.On("UpdateBook", bookId, bookRequestUpdate).Return(dtos.BookResponse{}, utils.ErrAuthorIdNotFound)

	bookControllerTest := bookController{bookService: bookServiceMock}

	request, err := http.NewRequest("PUT", fmt.Sprintf("/book/%v", bookId), bodyRequest)
	request.Header.Add("Content-type", "application/json")

	recorder := httptest.NewRecorder()
	e := echo.New()
	e.PUT("/book/:id", bookControllerTest.UpdateBook)
	e.ServeHTTP(recorder, request)

	require.NoError(t, err)
	require.Equal(t, http.StatusBadRequest, recorder.Code)

}
