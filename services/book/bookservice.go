package services

import (
	"errors"
	"fmt"
	"github/brunojoenk/golang-test/models/dtos"
	"github/brunojoenk/golang-test/models/entities"
	authorrepo "github/brunojoenk/golang-test/repository/author"
	bookrepo "github/brunojoenk/golang-test/repository/book"
	"github/brunojoenk/golang-test/utils"

	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type IBookService interface {
	CreateBook(bookRequestCreate dtos.BookRequestCreate) (dtos.BookResponse, error)
	GetAllBooks(filter dtos.GetBooksFilter) (dtos.BookResponseMetadata, error)
	DeleteBook(id int) error
	GetBook(id int) (dtos.BookResponse, error)
	UpdateBook(id int, bookRequestUpdate dtos.BookRequestUpdate) (dtos.BookResponse, error)
}

type bookService struct {
	authorDb authorrepo.IAuthorRepository
	bookDb   bookrepo.IBookRepository
}

// NewBookService Service Constructor
func NewBookService(db *gorm.DB) IBookService {
	authorRepo := authorrepo.NewAuthorRepository(db)
	bookRepo := bookrepo.NewBookRepository(db)
	return &bookService{
		authorDb: authorRepo,
		bookDb:   bookRepo,
	}
}

func (b *bookService) CreateBook(bookRequestCreate dtos.BookRequestCreate) (dtos.BookResponse, error) {
	var authors []entities.Author
	for _, authorId := range bookRequestCreate.Authors {
		author, err := b.authorDb.GetAuthor(authorId)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return dtos.BookResponse{}, utils.ErrAuthorIdNotFound
			}
			log.Error("Error on get author from repo: ", err.Error())
			return dtos.BookResponse{}, err
		}
		authors = append(authors, author)
	}

	book := entities.Book{
		Name:            bookRequestCreate.Name,
		Edition:         bookRequestCreate.Edition,
		PublicationYear: bookRequestCreate.PublicationYear,
		Authors:         authors,
	}

	createdBook, err := b.bookDb.CreateBook(book)
	if err != nil {
		log.Error("Error on create book from repo: ", err.Error())
		return dtos.BookResponse{}, err
	}

	dtosBookResponse := dtos.BookResponse{
		Id:              createdBook.Id,
		Name:            createdBook.Name,
		Edition:         createdBook.Edition,
		PublicationYear: createdBook.PublicationYear,
	}

	return dtosBookResponse, nil
}

func (b *bookService) GetAllBooks(filter dtos.GetBooksFilter) (dtos.BookResponseMetadata, error) {

	filter.Pagination.ValidValuesAndSetDefault()
	books, err := b.bookDb.GetAllBooks(filter)
	if err != nil {
		log.Error("Error on get all books from repo: ", err.Error())
		return dtos.BookResponseMetadata{}, err
	}

	booksResponse := make([]dtos.BookResponse, len(books))
	for i, book := range books {

		var authors string
		for i, author := range book.Authors {
			if i == 0 {
				authors = author.Name
				continue
			}
			authors += fmt.Sprintf(" | %s", author.Name)
		}

		bookResponse := dtos.BookResponse{
			Id:              book.Id,
			Name:            book.Name,
			Edition:         book.Edition,
			PublicationYear: book.PublicationYear,
			Authors:         authors,
		}

		booksResponse[i] = bookResponse
	}

	booksResponseMetadata := dtos.BookResponseMetadata{
		Books:      booksResponse,
		Pagination: filter.Pagination,
	}

	return booksResponseMetadata, nil
}

func (b *bookService) DeleteBook(id int) error {
	return b.bookDb.DeleteBook(id)
}

func (b *bookService) GetBook(id int) (dtos.BookResponse, error) {
	book, err := b.bookDb.GetBook(id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return dtos.BookResponse{}, utils.ErrBookIdNotFound
		}
		log.Error("Error on get book from repo: ", err.Error())
		return dtos.BookResponse{}, err
	}

	var authors string
	for i, author := range book.Authors {
		if i == 0 {
			authors = author.Name
			continue
		}
		authors += fmt.Sprintf(" | %s", author.Name)
	}

	bookResponse := dtos.BookResponse{
		Id:              book.Id,
		Name:            book.Name,
		Edition:         book.Edition,
		PublicationYear: book.PublicationYear,
		Authors:         authors,
	}

	return bookResponse, nil
}

func (b *bookService) UpdateBook(id int, bookRequestUpdate dtos.BookRequestUpdate) (dtos.BookResponse, error) {
	book, err := b.bookDb.GetBook(id)

	if err != nil {
		log.Error("Error on get book from repo: ", err.Error())
		return dtos.BookResponse{}, err
	}

	var authors []entities.Author
	for _, authorId := range bookRequestUpdate.Authors {
		author, err := b.authorDb.GetAuthor(authorId)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return dtos.BookResponse{}, utils.ErrAuthorIdNotFound
			}
			log.Error("Error on get author from repo: ", err.Error())
			return dtos.BookResponse{}, err
		}
		authors = append(authors, author)
	}

	book.Name = bookRequestUpdate.Name
	book.Edition = bookRequestUpdate.Edition
	book.PublicationYear = bookRequestUpdate.PublicationYear

	updatedBook, err := b.bookDb.UpdateBook(book, authors)

	if err != nil {
		log.Error("Error on update book from repo: ", err.Error())
		return dtos.BookResponse{}, err
	}

	var authorsString string
	for i, author := range updatedBook.Authors {
		if i == 0 {
			authorsString = author.Name
			continue
		}
		authorsString += fmt.Sprintf(" | %s", author.Name)
	}

	dtosBookResponse := dtos.BookResponse{
		Id:              updatedBook.Id,
		Name:            updatedBook.Name,
		Edition:         updatedBook.Edition,
		PublicationYear: updatedBook.PublicationYear,
		Authors:         authorsString,
	}

	return dtosBookResponse, nil

}
