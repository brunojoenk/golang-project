package services

import (
	"fmt"
	"github/brunojoenk/golang-test/models/dtos"
	"github/brunojoenk/golang-test/models/entities"
	authorrepo "github/brunojoenk/golang-test/repository/author"
	bookrepo "github/brunojoenk/golang-test/repository/book"
	"github/brunojoenk/golang-test/utils"

	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type GetAuthor func(id int) (*entities.Author, error)
type CreateBook func(book *entities.Book) error
type GetAllBooks func(filter dtos.GetBooksFilter) ([]entities.Book, error)
type DeleteBook func(id int) error
type GetBook func(id int) (*entities.Book, error)
type UpdateBook func(book *entities.Book, authors []*entities.Author) error

type BookService struct {
	getAuthorRepo   GetAuthor
	createBookRepo  CreateBook
	getAllBooksRepo GetAllBooks
	deleteBookRepo  DeleteBook
	getBookRepo     GetBook
	updateBookRepo  UpdateBook
}

// NewBookService Service Constructor
func NewBookService(db *gorm.DB) *BookService {
	bookRepo := bookrepo.NewBookRepository(db)
	authorRepo := authorrepo.NewAuthorRepository(db)
	return &BookService{
		getAuthorRepo:   authorRepo.GetAuthor,
		createBookRepo:  bookRepo.CreateBook,
		getAllBooksRepo: bookRepo.GetAllBooks,
		deleteBookRepo:  bookRepo.DeleteBook,
		getBookRepo:     bookRepo.GetBook,
		updateBookRepo:  bookRepo.UpdateBook,
	}
}

func (b *BookService) CreateBook(bookRequestCreate dtos.BookRequestCreateUpdate) error {
	var authors []*entities.Author
	for _, authorId := range bookRequestCreate.Authors {
		author, err := b.getAuthorRepo(authorId)
		if err != nil {
			log.Error("Error on get author from repo: ", err.Error())
			return err
		}
		if author.Id == 0 {
			return utils.ErrAuthorIdNotFound
		}
		authors = append(authors, author)
	}

	// Create book.
	book := entities.Book{
		Name:            bookRequestCreate.Name,
		Edition:         bookRequestCreate.Edition,
		PublicationYear: bookRequestCreate.PublicationYear,
		Authors:         authors,
	}

	return b.createBookRepo(&book)
}

func (b *BookService) GetAllBooks(filter dtos.GetBooksFilter) (*dtos.BookResponseMetadata, error) {

	filter.Pagination.ValidValuesAndSetDefault()
	books, err := b.getAllBooksRepo(filter)
	if err != nil {
		log.Error("Error on get all books from repo: ", err.Error())
		return nil, err
	}

	booksResponse := make([]dtos.BookResponse, 0)
	for _, book := range books {

		var authors string
		for i, author := range book.Authors {
			if i == 0 {
				authors = author.Name
				continue
			}
			authors += fmt.Sprintf(" | %s", author.Name)
		}

		bookResponse := &dtos.BookResponse{
			Name:            book.Name,
			Edition:         book.Edition,
			PublicationYear: book.PublicationYear,
			Authors:         authors,
		}

		booksResponse = append(booksResponse, *bookResponse)
	}

	booksResponseMetadata := &dtos.BookResponseMetadata{
		Books:      booksResponse,
		Pagination: filter.Pagination,
	}

	return booksResponseMetadata, nil
}

func (b *BookService) DeleteBook(id int) error {
	return b.deleteBookRepo(id)
}

func (b *BookService) GetBook(id int) (*dtos.BookResponse, error) {
	book, err := b.getBookRepo(id)

	if err != nil {
		log.Error("Error on get book from repo: ", err.Error())
		return nil, err
	}

	if book.Id == 0 {
		return nil, utils.ErrBookIdNotFound
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
		Name:            book.Name,
		Edition:         book.Edition,
		PublicationYear: book.PublicationYear,
		Authors:         authors,
	}

	return &bookResponse, nil
}

func (b *BookService) UpdateBook(id int, bookRequestUpdate dtos.BookRequestCreateUpdate) error {
	book, err := b.getBookRepo(id)

	if err != nil {
		log.Error("Error on get book from repo: ", err.Error())
		return err
	}

	var authors []*entities.Author
	for _, authorId := range bookRequestUpdate.Authors {
		author, err := b.getAuthorRepo(authorId)
		if err != nil {
			log.Error("Error on get author from repo: ", err.Error())
			return err
		}
		if author.Id == 0 {
			return utils.ErrAuthorIdNotFound
		}
		authors = append(authors, author)
	}

	book.Name = bookRequestUpdate.Name
	book.Edition = bookRequestUpdate.Edition
	book.PublicationYear = bookRequestUpdate.PublicationYear

	err = b.updateBookRepo(book, authors)

	if err != nil {
		log.Error("Error on update book from repo: ", err.Error())
		return err
	}

	return nil
}
