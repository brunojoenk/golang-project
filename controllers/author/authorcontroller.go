package controllers

import (
	"fmt"
	"github/brunojoenk/golang-test/models/dtos"
	authorservice "github/brunojoenk/golang-test/services/author"
	"net/http"
	"os"

	_ "github/brunojoenk/golang-test/docs"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type IAuthorController interface {
	GetAllAuthors(c echo.Context) error
	ReadCsvHandler(c echo.Context) error
}

type authorController struct {
	authorService authorservice.IAuthorService
}

// NewAuthorController Controller Constructor
func NewAuthorController(d *gorm.DB) *authorController {
	return &authorController{authorService: authorservice.NewAuthorService(d)}
}

// GetAllAuthors godoc
// @Summary Show all the authors with paginations.
// @Description Show all the authors with paginations.
// @Tags Authors
// @Accept */*
// @Produce json
// @Param   name     query     string     false  "search authors by name"     example(string)
// @Param   page     query     int     false  "page list"     example(1) minimum(1)
// @Param   limit     query     int     false  "page size"     example(1) minimum(1)
// @Success 200 {object} dtos.AuthorResponseMetadata
// @Failure 400 {object} string
// @Failure 500 {object} string
// @Router /authors [get]
func (a *authorController) GetAllAuthors(c echo.Context) error {

	var filter dtos.GetAuthorsFilter
	err := c.Bind(&filter)
	if err != nil {
		c.Logger().Warn("Error on bind query to filter author: %s", err.Error())
		return c.JSON(http.StatusBadRequest, "Invalid parameter")
	}

	authorsResponse, err := a.authorService.GetAllAuthors(filter)
	if err != nil {
		c.Logger().Error("Error get all author: %s", err.Error())
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, authorsResponse)
}

// Import authors from authors.csv godoc
// @Summary Import authors from authors.csv.
// @Description Import authors from authors.csv.
// @Tags Authors
// @Accept */*
// @Produce json
// @Success 200 {array} dtos.AuthorResponseMetadata
// @Router /authors/import [post]
func (a *authorController) ReadCsvHandler(c echo.Context) error {

	authorsFilePath := os.Getenv("AUTHORS_FILE_PATH")
	if authorsFilePath == "" {
		c.Logger().Info("Setting default value for author file path (env AUTHORS_FILE_PATH)")
		//Set default, safe mode. When run locally, this env is exported on makefile
		authorsFilePath = "./data/authorsreduced.csv"
	}
	totalAuthorsAdded, err := a.authorService.ImportAuthorsFromCSVFile(authorsFilePath)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, fmt.Sprintf("Error on import authors: %s", err.Error()))
	}
	return c.JSON(http.StatusOK, dtos.AuthorImportResponse{
		Msg:   "Authors imported",
		Total: totalAuthorsAdded,
	})
}
