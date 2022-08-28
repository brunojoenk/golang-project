package services

import (
	"fmt"
	"github/brunojoenk/golang-test/models"
	authorrepo "github/brunojoenk/golang-test/repository/author"
	bookrepo "github/brunojoenk/golang-test/repository/book"
	"github/brunojoenk/golang-test/utils"

	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type GetAuthor func(id int) (*models.Author, error)
type CreateBook func(book *models.Book) error
type GetAllBooks func(filter models.GetBooksFilter) ([]models.Book, error)
type DeleteBook func(id int) error
type GetBook func(id int) (*models.Book, error)
type UpdateBook func(book *models.Book, authors []*models.Author) error

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

func (b *BookService) CreateBook(bookRequestCreate models.BookRequestCreateUpdate) error {
	var authors []*models.Author
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
	book := models.Book{
		Name:            bookRequestCreate.Name,
		Edition:         bookRequestCreate.Edition,
		PublicationYear: bookRequestCreate.PublicationYear,
		Authors:         authors,
	}

	return b.createBookRepo(&book)
}

func (b *BookService) GetAllBooks(filter models.GetBooksFilter) (*models.BookResponseMetadata, error) {

	filter.Pagination.ValidValuesAndSetDefault()
	books, err := b.getAllBooksRepo(filter)
	if err != nil {
		log.Error("Error on get all books from repo: ", err.Error())
		return nil, err
	}

	booksResponse := make([]models.BookResponse, 0)
	for _, book := range books {

		var authors string
		for i, author := range book.Authors {
			if i == 0 {
				authors = author.Name
				continue
			}
			authors += fmt.Sprintf(" | %s", author.Name)
		}

		bookResponse := &models.BookResponse{
			Name:            book.Name,
			Edition:         book.Edition,
			PublicationYear: book.PublicationYear,
			Authors:         authors,
		}

		booksResponse = append(booksResponse, *bookResponse)
	}

	booksResponseMetadata := &models.BookResponseMetadata{
		Books:      booksResponse,
		Pagination: filter.Pagination,
	}

	return booksResponseMetadata, nil
}

func (b *BookService) DeleteBook(id int) error {
	return b.deleteBookRepo(id)
}

func (b *BookService) GetBook(id int) (*models.BookResponse, error) {
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

	bookResponse := models.BookResponse{
		Name:            book.Name,
		Edition:         book.Edition,
		PublicationYear: book.PublicationYear,
		Authors:         authors,
	}

	return &bookResponse, nil
}

func (b *BookService) UpdateBook(id int, bookRequestUpdate models.BookRequestCreateUpdate) error {
	book, err := b.getBookRepo(id)

	if err != nil {
		log.Error("Error on get book from repo: ", err.Error())
		return err
	}

	var authors []*models.Author
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
