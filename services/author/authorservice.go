package services

import (
	"encoding/csv"
	"github/brunojoenk/golang-test/models/dtos"
	"github/brunojoenk/golang-test/models/entities"
	authorrepo "github/brunojoenk/golang-test/repository/author"
	"os"
	"sync"

	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

var BATCH_SIZE = 2000

type GetAllAuthors func(filter dtos.GetAuthorsFilter) ([]entities.Author, error)
type CreateAuthorInBatch func(author []entities.Author, batchSize int) error

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

func (a *AuthorService) ImportAuthorsFromCSVFile(file string) (int, error) {

	f, err := os.Open(file)

	if err != nil {
		log.Error("Error on open file: ", err.Error())
		return 0, err
	}

	defer f.Close()

	fcsv := csv.NewReader(f)
	fcsv.Comma = ';'

	numWorkers := 20
	jobs := make(chan []entities.Author, numWorkers)
	res := make(chan []entities.Author)

	worker := func(jobs <-chan []entities.Author, results chan<- []entities.Author) error {
		for job := range jobs {
			err := a.createAuthorInBatchRepo(job, len(job))
			if err != nil {
				log.Error("Error on create author in batch repository: ", err.Error())
				return err
			}
			results <- job
		}
		return nil
	}

	var errOnBatch error
	var wg sync.WaitGroup
	// init workers
	for w := 0; w < numWorkers; w++ {
		wg.Add(1)
		go func() {
			// this line will exec when chan `res` processed
			defer wg.Done()
			errOnBatch = worker(jobs, res)
		}()
	}

	go func() {
		authorsAddedMap := make(map[string]bool, 0)
		records, err := fcsv.ReadAll()
		if err != nil {
			log.Error("Error on read all csv: ", err.Error())
			return
		}
		for _, record := range records {
			a.processRecord(record, authorsAddedMap, jobs)
		}
		close(jobs) // close jobs to signal workers that no more job are incoming.
	}()

	go func() {
		wg.Wait()
		close(res) // when you close(res) it breaks the below loop.
	}()

	authors := make([]entities.Author, 0)
	for r := range res {
		authors = append(authors, r...)
	}

	return len(authors), errOnBatch
}

func (a *AuthorService) processRecord(record []string, authorsAddedMap map[string]bool, jobs chan []entities.Author) {
	batchToCreate := make([]entities.Author, 0)
	for index, name := range record {
		if a.isAuthorNotAdded(authorsAddedMap, name) {
			authorsAddedMap[name] = true
			batchToCreate = append(batchToCreate, entities.Author{Name: name})
		}
		if a.canCreateInBatch(index, len(record), len(batchToCreate)) {
			jobs <- batchToCreate
			batchToCreate = make([]entities.Author, 0)
		}
	}
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
