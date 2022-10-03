package handlers

import (
	authorcontroller "github/brunojoenk/golang-test/controllers/author"
	bookcontroller "github/brunojoenk/golang-test/controllers/book"

	_ "github/brunojoenk/golang-test/docs"

	// echo-swagger middleware

	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
	"gorm.io/gorm"
)

type Handler struct {
	authorController authorcontroller.IAuthorController
	bookController   bookcontroller.IBookController
}

func New(db *gorm.DB) *Handler {
	return &Handler{authorController: authorcontroller.NewAuthorController(db), bookController: bookcontroller.NewBookController(db)}
}

func (h *Handler) HandleControllers(e *echo.Echo) {
	e.POST("/authors/import", h.authorController.ReadCsvHandler)
	e.GET("/authors", h.authorController.GetAllAuthors)

	e.POST("/book", h.bookController.CreateBook)
	e.GET("/books", h.bookController.GetAllBooks)
	e.GET("/book/:id", h.bookController.GetBook)
	e.PUT("/book/:id", h.bookController.UpdateBook)
	e.DELETE("/book/:id", h.bookController.DeleteBook)

	e.GET("/swagger/*", echoSwagger.WrapHandler)
}
