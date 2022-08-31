package controllers

import (
	"fmt"
	"github/brunojoenk/golang-test/models/dtos"
	bookservice "github/brunojoenk/golang-test/services/book"
	"github/brunojoenk/golang-test/utils"
	"net/http"
	"strconv"

	"github.com/pkg/errors"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type CreateBook func(bookRequestCreate dtos.BookRequestCreateUpdate) error
type GetAllBooks func(filter dtos.GetBooksFilter) (*dtos.BookResponseMetadata, error)
type DeleteBook func(id int) error
type GetBook func(id int) (*dtos.BookResponse, error)
type UpdateBook func(id int, bookRequestUpdate dtos.BookRequestCreateUpdate) error

type BookController struct {
	createBookRepo  CreateBook
	getAllBooksRepo GetAllBooks
	deleteBookRepo  DeleteBook
	getBookRepo     GetBook
	updateBookRepo  UpdateBook
}

// NewBookController Controller Constructor
func NewBookController(db *gorm.DB) *BookController {
	repo := bookservice.NewBookService(db)
	return &BookController{createBookRepo: repo.CreateBook,
		getAllBooksRepo: repo.GetAllBooks,
		deleteBookRepo:  repo.DeleteBook,
		getBookRepo:     repo.GetBook,
		updateBookRepo:  repo.UpdateBook}
}

// CreateBook godoc
// @Summary Create a book.
// @Description Create a book.
// @Tags Books
// @Accept json
// @Produce json
// @Param request body dtos.BookRequestCreateUpdate true "query params"
// @Success 201 {object} string
// @Failure 400 {object} string
// @Failure 500 {object} string
// @Router /book [post]
func (b *BookController) CreateBook(c echo.Context) error {

	bookRequestCreate := new(dtos.BookRequestCreateUpdate)
	if err := c.Bind(bookRequestCreate); err != nil {
		c.Logger().Warn("Error on bind body to create book: %s", err.Error())
		return c.JSON(http.StatusBadRequest, fmt.Sprintf("Request body to crate a book is invalid: %s", err.Error()))
	}

	err := b.createBookRepo(*bookRequestCreate)

	if err != nil {
		if errors.Is(err, utils.ErrAuthorIdNotFound) {
			return c.JSON(http.StatusBadRequest, "Author id not found to create book with author")
		}
		c.Logger().Error("Error on create book: %s", err.Error())
		return c.JSON(http.StatusInternalServerError, "Error on create book. Please, contact system admin")
	}

	return c.JSON(http.StatusCreated, "Created")
}

// GetAllBooks godoc
// @Summary Show all the books with paginations.
// @Description Show all the books with paginations.
// @Tags Books
// @Accept */*
// @Produce json
// @Param   name     query     string     false  "search book by name"     example(string)
// @Param   edition     query     string     false  "search book by edition"     example(string)
// @Param   publication_year     query     int     false  "search book by publication year"     example(1) minimum(1)
// @Param   author     query     string     false  "search book by author"     example(string)
// @Param   page     query     int     false  "page list"     example(1) minimum(1)
// @Param   limit     query     int     false  "page size"     example(1) minimum(1)
// @Success 200 {object} dtos.BookResponseMetadata
// @Failure 400 {object} string
// @Failure 500 {object} string
// @Router /books [get]
func (b *BookController) GetAllBooks(c echo.Context) error {
	var filter dtos.GetBooksFilter
	err := c.Bind(&filter)
	if err != nil {
		c.Logger().Warn("Error on bind query to filter all books: %s", err.Error())
		return c.JSON(http.StatusBadRequest, fmt.Sprintf("Invalid query parameters: %s", err.Error()))
	}

	booksResponse, err := b.getAllBooksRepo(filter)

	if err != nil {
		c.Logger().Error("Error on get all books: %s", err.Error())
		return c.JSON(http.StatusInternalServerError, "Error on get all books. Please, contact admin")
	}

	return c.JSON(http.StatusOK, booksResponse)
}

// DeleteBook godoc
// @Summary Delete a book.
// @Description Delete a book.
// @Tags Books
// @Accept */*
// @Produce json
// @Param id   path int true "Book ID"
// @Success 200 {object} string
// @Failure 400 {object} string
// @Failure 500 {object} string
// @Router /book/{id} [delete]
func (b *BookController) DeleteBook(c echo.Context) error {

	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.Logger().Warn("Error on parse parameters id on delete book %s", err.Error())
		return c.JSON(http.StatusBadRequest, "Invalid query parameter id")
	}

	err = b.deleteBookRepo(id)

	if err != nil {
		c.Logger().Error("Error on delete book: %s", err.Error())
		return c.JSON(http.StatusInternalServerError, "Erro on delete book. Please, contact system admin")
	}

	return c.JSON(http.StatusOK, "Deleted")
}

// GetBook godoc
// @Summary Get a book.
// @Description gET a book.
// @Tags Books
// @Accept */*
// @Produce json
// @Param id   path int true "Book ID"
// @Success 200 {object} dtos.BookResponse
// @Failure 400 {object} string
// @Failure 500 {object} string
// @Router /book/{id} [get]
func (b *BookController) GetBook(c echo.Context) error {

	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.Logger().Warn("Error on parse parameters id on get book %s", err.Error())
		return c.JSON(http.StatusBadRequest, "Invalid query parameter id")
	}

	bookResponse, err := b.getBookRepo(id)

	if err != nil {
		if errors.Is(err, utils.ErrBookIdNotFound) {
			return c.JSON(http.StatusBadRequest, "Author ID not found")
		}
		c.Logger().Error("Error on get book %s", err.Error())
		return c.JSON(http.StatusInternalServerError, "Error on get book. Please contact system admin")
	}

	return c.JSON(http.StatusOK, bookResponse)
}

// UpdateBook godoc
// @Summary Update a book.
// @Description Update a book.
// @Tags Books
// @Accept */*
// @Produce json
// @Param id   path int true "Book ID"
// @Param request body dtos.BookRequestCreateUpdate true "query params"
// @Success 200 {object} string
// @Failure 400 {object} string
// @Failure 500 {object} string
// @Router /book/{id} [put]
func (b *BookController) UpdateBook(c echo.Context) error {

	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.Logger().Warn("Error on parse parameters id on update book %s", err.Error())
		return c.JSON(http.StatusBadRequest, "Invalid query parameter id")
	}

	bookRequestUpdate := new(dtos.BookRequestCreateUpdate)
	if err := c.Bind(bookRequestUpdate); err != nil {
		c.Logger().Warn("Error on parse body on update book %s", err.Error())
		return c.JSON(http.StatusBadRequest, fmt.Sprintf("Error on parse body to update book: %s", err.Error()))
	}

	err = b.updateBookRepo(id, *bookRequestUpdate)

	if err != nil {
		if errors.Is(err, utils.ErrAuthorIdNotFound) {
			return c.JSON(http.StatusBadRequest, err.Error())
		}
		c.Logger().Error("Error on update book %s", err.Error())
		return c.JSON(http.StatusInternalServerError, "Error on update book. Please contact system admin")
	}

	return c.JSON(http.StatusOK, "Updated")
}
