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

var BATCH_SIZE = 2000

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

	var names []string

	for _, record := range records {
		names, err = a.processRecord(record)
		if err != nil {
			return nil, err
		}
	}

	return names, err
}

func (a *AuthorService) processRecord(record []string) ([]string, error) {
	authorsAddedMap := make(map[string]bool, 0)
	batchToCreate := make([]*entities.Author, 0)
	for index, name := range record {
		if a.authorNotAdded(authorsAddedMap, name) {
			authorsAddedMap[name] = true
			batchToCreate = append(batchToCreate, &entities.Author{Name: name})
		}
		if a.canCreateInBatch(index, len(record)) {
			err := a.createAuthorInBatchRepo(batchToCreate, index)
			if err != nil {
				log.Error("Error on create author in batch repository: ", err.Error())
				return nil, err
			}
			batchToCreate = make([]*entities.Author, 0)
		}
	}

	namesAdded := make([]string, 0)
	for name := range authorsAddedMap {
		namesAdded = append(namesAdded, name)
	}

	return namesAdded, nil
}

func (a *AuthorService) canCreateInBatch(index, recordSize int) bool {
	return (index > 0 && a.isCounterEqualBatchSize(index)) || a.isLastItemToIterate(index, recordSize)
}

func (a *AuthorService) isCounterEqualBatchSize(index int) bool {
	return index%BATCH_SIZE == 0
}

func (a *AuthorService) isLastItemToIterate(index, recordSize int) bool {
	return index == (recordSize - 1)
}

func (a *AuthorService) authorNotAdded(authorsAddedMap map[string]bool, name string) bool {
	return !authorsAddedMap[name]
}
