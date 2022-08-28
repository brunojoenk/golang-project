package repository

import (
	"github/brunojoenk/golang-test/models"
	"strings"

	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// AuthorsRepository Author Repository
type AuthorRepository struct {
	db *gorm.DB
}

// NewAuthorsRepository Repository Constructor
func NewAuthorRepository(d *gorm.DB) *AuthorRepository {
	return &AuthorRepository{db: d}
}

func (a *AuthorRepository) CreateAuthorInBatch(author []*models.Author, batchSize int) error {

	if result := a.db.CreateInBatches(author, batchSize); result.Error != nil {
		log.Error("Error on create authors in batch: ", result.Error.Error())
		return result.Error
	}

	return nil
}

func (a *AuthorRepository) GetAuthor(id int) (*models.Author, error) {

	var author models.Author

	if result := a.db.Find(&author, id); result.Error != nil {
		log.Error("Error on get author: ", result.Error.Error())
		return nil, result.Error
	}

	return &author, nil
}

// Get authors
func (a *AuthorRepository) GetAllAuthors(filter models.GetAuthorsFilter) ([]models.Author, error) {

	var authors []models.Author
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
