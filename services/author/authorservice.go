package services

import (
	"encoding/csv"
	"github/brunojoenk/golang-test/models/dtos"
	"github/brunojoenk/golang-test/models/entities"
	authorrepo "github/brunojoenk/golang-test/repository/author"
	"os"

	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type GetAllAuthors func(filter dtos.GetAuthorsFilter) ([]entities.Author, error)
type CreateAuthorInBatch func(author []*entities.Author, batchSize int) error

type AuthorService struct {
	getAllAuthorsRepository GetAllAuthors
	createAuthorInBatchRepo CreateAuthorInBatch
}

// NewBookService Service Constructor
func NewAuthorService(db *gorm.DB) *AuthorService {
	repo := authorrepo.NewAuthorRepository(db)
	return &AuthorService{getAllAuthorsRepository: repo.GetAllAuthors, createAuthorInBatchRepo: repo.CreateAuthorInBatch}
}

func (a *AuthorService) GetAllAuthors(filter dtos.GetAuthorsFilter) (*dtos.AuthorResponseMetadata, error) {

	filter.Pagination.ValidValuesAndSetDefault()
	authors, err := a.getAllAuthorsRepository(filter)
	if err != nil {
		log.Error("Error on get all authors from repositoriy: ", err.Error())
		return nil, err
	}

	authorsResponse := make([]dtos.AuthorResponse, 0)
	for _, a := range authors {
		authorResponse := &dtos.AuthorResponse{
			Id:   a.Id,
			Name: a.Name,
		}
		authorsResponse = append(authorsResponse, *authorResponse)
	}

	authorResponseMetada := &dtos.AuthorResponseMetadata{
		Authors:    authorsResponse,
		Pagination: filter.Pagination,
	}

	return authorResponseMetada, nil
}

func (a *AuthorService) ImportAuthorsFromCSVFile(file string) ([]string, error) {

	f, err := os.Open(file)

	if err != nil {
		log.Error("Error on open file: ", err.Error())
		return nil, err
	}

	defer f.Close()

	r := csv.NewReader(f)
	r.Comma = ';'

	records, err := r.ReadAll()

	if err != nil {
		log.Error("Error on readall from file: ", err.Error())
		return nil, err
	}

	names := make([]string, 0)
	mapper := make(map[string]bool, 0)
	batchSize := 2000

	for _, record := range records {
		count := 0
		var batch []*entities.Author
		for i, name := range record {
			if !mapper[name] {
				mapper[name] = true
				count++
				batch = append(batch, &entities.Author{Name: name})
				names = append(names, name)
			}
			if count == batchSize || i == (len(record)-1) {
				err := a.createAuthorInBatchRepo(batch, count)
				if err != nil {
					log.Error("Error on create author in batch repository: ", err.Error())
					return nil, err
				}
				batch = make([]*entities.Author, 0)
				count = 0
			}

		}

	}

	return names, err
}
