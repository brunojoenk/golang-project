package repository

import (
	"github/brunojoenk/golang-test/models/dtos"
	"github/brunojoenk/golang-test/models/entities"
	"strings"

	log "github.com/sirupsen/logrus"

	"gorm.io/gorm"
)

// BookRepository Books Repository
type BookRepository struct {
	db *gorm.DB
}

// NewBooksRepository Repository Constructor
func NewBookRepository(db *gorm.DB) *BookRepository {
	return &BookRepository{db: db}
}

// Create book
func (b *BookRepository) CreateBook(book *entities.Book) error {

	if result := b.db.Create(&book); result.Error != nil {
		log.Error("Error on create book: ", result.Error.Error())
		return result.Error
	}

	return nil
}

func (b *BookRepository) UpdateBook(book *entities.Book, authors []*entities.Author) error {

	if err := b.db.Model(&book).Association("Authors").Clear(); err != nil {
		log.Error("Error on clear authors from book: ", err.Error())
		return err
	}

	book.Authors = authors

	if result := b.db.Save(&book); result.Error != nil {
		log.Error("Error on update book: ", result.Error.Error())
		return result.Error
	}

	return nil
}

func (b *BookRepository) GetBook(id int) (*entities.Book, error) {
	var book entities.Book

	if result := b.db.Preload("Authors").First(&book, id); result.Error != nil {
		log.Error("Error on preload authors from book: ", result.Error.Error())
		return nil, result.Error
	}

	return &book, nil
}

// Get books
func (b *BookRepository) GetAllBooks(filter dtos.GetBooksFilter) ([]entities.Book, error) {

	var books []entities.Book
	toExec := b.db

	if strings.TrimSpace(filter.Author) != "" {
		toExec = toExec.Joins(
			"JOIN author_book ON author_book.book_id = books.id " +
				"JOIN authors ON authors.id = author_book.author_id")
		toExec = toExec.Where("LOWER(authors.name) LIKE ?", "%"+strings.ToLower(filter.Author)+"%")
	}

	if strings.TrimSpace(filter.Name) != "" {
		toExec = toExec.Where("LOWER(books.name) LIKE ?", "%"+strings.ToLower(filter.Name)+"%")
	}

	if strings.TrimSpace(filter.Edition) != "" {
		toExec = toExec.Where("LOWER(books.edition) LIKE ?", "%"+strings.ToLower(filter.Edition)+"%")
	}

	if filter.PublicationYear > 0 {
		toExec = toExec.Where("books.publication_year = ?", filter.PublicationYear)
	}

	toExec = toExec.Offset((filter.Page - 1) * filter.Limit).Limit(filter.Limit).Order("name asc")

	if result := toExec.Preload("Authors").Find(&books); result.Error != nil {
		log.Error("Error on preload authors from book (filter): ", result.Error.Error())
		return nil, result.Error
	}

	return books, nil
}

func (b *BookRepository) DeleteBook(id int) error {
	var book entities.Book

	if result := b.db.First(&book, id); result.Error != nil {
		log.Error("Error on get book to delete: ", result.Error.Error())
		return result.Error
	}

	if result := b.db.Exec("DELETE FROM author_book WHERE author_book.book_id = $1", id); result.Error != nil {
		log.Error("Error on delete relations from author_book: ", result.Error.Error())
		return result.Error
	}

	if result := b.db.Delete(&book); result.Error != nil {
		log.Error("Error delete book: ", result.Error.Error())
		return result.Error
	}

	return nil
}
