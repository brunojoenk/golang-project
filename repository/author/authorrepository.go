package repository

import (
	"github/brunojoenk/golang-test/models/dtos"
	"github/brunojoenk/golang-test/models/entities"
	"strings"

	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type IAuthorRepository interface {
	CreateAuthorInBatch(author []entities.Author, batchSize int) error
	GetAuthor(id int) (entities.Author, error)
	GetAllAuthors(filter dtos.GetAuthorsFilter) ([]entities.Author, error)
}

// AuthorsRepository Author Repository
type AuthorRepository struct {
	db *gorm.DB
}

// NewAuthorsRepository Repository Constructor
func NewAuthorRepository(d *gorm.DB) IAuthorRepository {
	return &AuthorRepository{db: d}
}

func (a *AuthorRepository) CreateAuthorInBatch(author []entities.Author, batchSize int) error {

	if result := a.db.Clauses(clause.OnConflict{
		DoNothing: true,
		Columns:   []clause.Column{{Name: "name"}}}).
		CreateInBatches(author, batchSize); result.Error != nil {
		log.Error("Error on create authors in batch: ", result.Error.Error())
		return result.Error
	}

	return nil
}

func (a *AuthorRepository) GetAuthor(id int) (entities.Author, error) {

	var author entities.Author

	if result := a.db.Find(&author, id); result.Error != nil {
		log.Error("Error on get author: ", result.Error.Error())
		return author, result.Error
	}

	return author, nil
}

func (a *AuthorRepository) GetAllAuthors(filter dtos.GetAuthorsFilter) ([]entities.Author, error) {

	var authors []entities.Author
	toExec := a.db

	if strings.TrimSpace(filter.Name) != "" {
		toExec = toExec.Where("LOWER(name) LIKE ?", "%"+strings.ToLower(filter.Name)+"%")
	}

	toExec = toExec.Offset((filter.Page - 1) * filter.Limit).Limit(filter.Limit).Order("name asc")

	if result := toExec.Find(&authors); result.Error != nil {
		log.Error("Error on get all authors: ", result.Error.Error())
		return nil, result.Error
	}

	return authors, nil
}
