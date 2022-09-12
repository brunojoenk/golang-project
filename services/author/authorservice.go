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

type AuthorService struct {
	authorDb authorrepo.IAuthorRepository
}

// NewBookService Service Constructor
func NewAuthorService(db *gorm.DB) *AuthorService {
	repo := authorrepo.NewAuthorRepository(db)
	return &AuthorService{authorDb: repo}
}

func (a *AuthorService) GetAllAuthors(filter dtos.GetAuthorsFilter) (*dtos.AuthorResponseMetadata, error) {

	filter.Pagination.ValidValuesAndSetDefault()
	authors, err := a.authorDb.GetAllAuthors(filter)
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

func (a *AuthorService) ImportAuthorsFromCSVFile(file string) (int, error) {

	f, err := os.Open(file)

	if err != nil {
		log.Error("Error on open file: ", err.Error())
		return 0, err
	}

	defer f.Close()

	r := csv.NewReader(f)
	r.Comma = ';'

	records, err := r.ReadAll()

	if err != nil {
		log.Error("Error on readall from file: ", err.Error())
		return 0, err
	}

	authorsAddedMap := make(map[string]bool, 0)
	totalAuthors := 0
	for _, record := range records {
		numAdded, err := a.processRecord(record, authorsAddedMap)
		totalAuthors = totalAuthors + numAdded
		if err != nil {
			return 0, err
		}
	}

	return totalAuthors, err
}

func (a *AuthorService) processRecord(record []string, authorsAddedMap map[string]bool) (int, error) {
	totalAuthors := 0
	batchToCreate := make([]entities.Author, 0)
	for index, name := range record {
		if a.isAuthorNotAdded(authorsAddedMap, name) {
			authorsAddedMap[name] = true
			batchToCreate = append(batchToCreate, entities.Author{Name: name})
		}
		if a.canCreateInBatch(index, len(record), len(batchToCreate)) {
			err := a.authorDb.CreateAuthorInBatch(batchToCreate, len(batchToCreate))
			if err != nil {
				log.Error("Error on create author in batch repository: ", err.Error())
				return totalAuthors, err
			}
			totalAuthors = totalAuthors + (len(batchToCreate))
			batchToCreate = make([]entities.Author, 0)
		}
	}

	return totalAuthors, nil
}

func (a *AuthorService) canCreateInBatch(index, recordSize, batchSize int) bool {
	return a.isCounterEqualBatchSize(batchSize) || a.isLastItemToProcess(index, recordSize)
}

func (a *AuthorService) isCounterEqualBatchSize(batchSize int) bool {
	return batchSize > 0 && batchSize%BATCH_SIZE == 0
}

func (a *AuthorService) isLastItemToProcess(index, recordSize int) bool {
	return index == (recordSize - 1)
}

func (a *AuthorService) isAuthorNotAdded(authorsAddedMap map[string]bool, name string) bool {
	return !authorsAddedMap[name]
}
