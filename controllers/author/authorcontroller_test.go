package controllers

import (
	"errors"
	"fmt"
	"github/brunojoenk/golang-test/models/dtos"

	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	authorservicemock "github/brunojoenk/golang-test/services/author/mock"

	"encoding/json"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
)

func TestGetAllAuthors(t *testing.T) {
	var (
		authorId   = 8
		authorName = "bruno"
		authors    = dtos.AuthorResponseMetadata{Authors: []dtos.AuthorResponse{{Id: authorId, Name: authorName}}, Pagination: dtos.Pagination{Page: 1, Limit: 10}}
	)

	authorServiceMock := new(authorservicemock.AuthorServiceyMock)
	authorServiceMock.On("GetAllAuthors", dtos.GetAuthorsFilter{}).Return(authors, nil)

	authorControllerTest := authorController{authorService: authorServiceMock}

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := authorControllerTest.GetAllAuthors(c)

	require.NoError(t, err)
	require.Equal(t, http.StatusOK, rec.Code)

	respExpected, _ := json.Marshal(authors)
	require.Equal(t, fmt.Sprintf("%s%s", respExpected, "\n"), rec.Body.String())
}

func TestGetAllAuthorsErrorOnFilter(t *testing.T) {
	request, err := http.NewRequest("GET", "/authors?limit='a'", nil)
	recorder := httptest.NewRecorder()
	e := echo.New()
	authorControllerTest := authorController{}
	e.GET("/authors", authorControllerTest.GetAllAuthors)
	e.ServeHTTP(recorder, request)

	require.Equal(t, http.StatusBadRequest, recorder.Code)

	if err != nil {
		t.Errorf("Expected nil, received %s", err.Error())
	}

}

func TestGetAllAuthorsErrorOnService(t *testing.T) {
	errExpected := errors.New("error occurred")
	authorServiceMock := new(authorservicemock.AuthorServiceyMock)
	authorServiceMock.On("GetAllAuthors", dtos.GetAuthorsFilter{}).Return(dtos.AuthorResponseMetadata{}, errExpected)

	authorControllerTest := authorController{authorService: authorServiceMock}

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	authorControllerTest.GetAllAuthors(c)

	//require.ErrorIs(t, errExpected, err)
	require.Equal(t, http.StatusInternalServerError, rec.Code)
}

func TestImportReadCsvHandler(t *testing.T) {
	os.Setenv("AUTHORS_FILE_PATH", "./data/authorsreduced.csv")

	authorServiceMock := new(authorservicemock.AuthorServiceyMock)
	authorServiceMock.On("ImportAuthorsFromCSVFile", "./data/authorsreduced.csv").Return(6, nil)

	authorControllerTest := authorController{authorService: authorServiceMock}

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	authorControllerTest.ReadCsvHandler(c)

	time.Sleep(time.Millisecond * 100) //added to receive value of increment becase import is call by goroutine

	//require.ErrorIs(t, errExpected, err)
	require.Equal(t, http.StatusOK, rec.Code)
}

func TestImportReadCsvHandlerError(t *testing.T) {
	os.Setenv("AUTHORS_FILE_PATH", "")
	authorServiceMock := new(authorservicemock.AuthorServiceyMock)
	authorServiceMock.On("ImportAuthorsFromCSVFile", "./data/authorsreduced.csv").Return(0, errors.New("Error occurred"))

	authorControllerTest := authorController{authorService: authorServiceMock}

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	authorControllerTest.ReadCsvHandler(c)

	//require.ErrorIs(t, errExpected, err)
	require.Equal(t, http.StatusInternalServerError, rec.Code)
}
