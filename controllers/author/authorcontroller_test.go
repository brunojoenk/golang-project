package controllers

import (
	"errors"
	"fmt"
	"github/brunojoenk/golang-test/models"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"encoding/json"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var authorControllerTest *AuthorController

func init() {
	db, _, _ := sqlmock.New()

	dialector := postgres.New(postgres.Config{
		DSN:                  "sqlmock_db_0",
		DriverName:           "postgres",
		Conn:                 db,
		PreferSimpleProtocol: true,
	})
	gormDB, _ := gorm.Open(dialector, &gorm.Config{})

	authorControllerTest = NewAuthorController(gormDB)
}

func TestGetAllAuthors(t *testing.T) {
	var (
		authorId   = 8
		authorName = "bruno"
		authors    = models.AuthorResponseMetadata{Authors: []models.AuthorResponse{{Id: authorId, Name: authorName}}, Pagination: models.Pagination{Page: 1, Limit: 10}}
	)
	authorControllerTest.getAllAuthorsRepo = func(filter models.GetAuthorsFilter) (*models.AuthorResponseMetadata, error) {
		return &authors, nil
	}

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
	e.GET("/authors", authorControllerTest.GetAllAuthors)
	e.ServeHTTP(recorder, request)

	require.Equal(t, http.StatusBadRequest, recorder.Code)

	if err != nil {
		t.Errorf("Expected nil, received %s", err.Error())
	}

}

func TestGetAllAuthorsErrorOnService(t *testing.T) {
	errExpected := errors.New("error occurred")
	authorControllerTest.getAllAuthorsRepo = func(filter models.GetAuthorsFilter) (*models.AuthorResponseMetadata, error) {
		return nil, errExpected
	}

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	authorControllerTest.GetAllAuthors(c)

	//require.ErrorIs(t, errExpected, err)
	require.Equal(t, http.StatusInternalServerError, rec.Code)
}

func TestImportReadCsvHandler(t *testing.T) {
	expectedFuncCalled := 0
	authorControllerTest.importAuthorsFromCSVFile = func(file string) ([]string, error) {
		expectedFuncCalled = expectedFuncCalled + 1
		return nil, nil
	}

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	authorControllerTest.ReadCsvHandler(c)

	time.Sleep(time.Millisecond * 100) //added to receive value of increment becase import is call by goroutine

	//require.ErrorIs(t, errExpected, err)
	require.Equal(t, http.StatusOK, rec.Code)
	require.Equal(t, expectedFuncCalled, 1)
}

func TestImportReadCsvHandlerError(t *testing.T) {
	os.Setenv("AUTHORS_FILE_PATH", "")
	expectedFuncCalled := 0
	fileCalled := ""
	authorControllerTest.importAuthorsFromCSVFile = func(file string) ([]string, error) {
		expectedFuncCalled = expectedFuncCalled + 1
		fileCalled = file
		return nil, errors.New("Error occurred")
	}

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	authorControllerTest.ReadCsvHandler(c)

	//require.ErrorIs(t, errExpected, err)
	require.Equal(t, http.StatusInternalServerError, rec.Code)
	require.Equal(t, expectedFuncCalled, 1)
	require.Equal(t, "./data/authorsreduced.csv", fileCalled)
}
