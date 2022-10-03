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

type IAuthorService interface {
	GetAllAuthors(filter dtos.GetAuthorsFilter) (*dtos.AuthorResponseMetadata, error)
	ImportAuthorsFromCSVFile(file string) (int, error)
}

type authorService struct {
	authorDb authorrepo.IAuthorRepository
}

// NewBookService Service Constructor
func NewAuthorService(db *gorm.DB) IAuthorService {
	return &authorService{authorDb: authorrepo.NewAuthorRepository(db)}
}

func (a *authorService) GetAllAuthors(filter dtos.GetAuthorsFilter) (*dtos.AuthorResponseMetadata, error) {

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

func (a *authorService) ImportAuthorsFromCSVFile(file string) (int, error) {

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
			err := a.authorDb.CreateAuthorInBatch(job, len(job))
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

func (a *authorService) processRecord(record []string, authorsAddedMap map[string]bool, jobs chan []entities.Author) {
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

func (a *authorService) canCreateInBatch(index, recordSize, batchSize int) bool {
	return a.isCounterEqualBatchSize(batchSize) || a.isLastItemToProcess(index, recordSize)
}

func (a *authorService) isCounterEqualBatchSize(batchSize int) bool {
	return batchSize > 0 && batchSize%BATCH_SIZE == 0
}

func (a *authorService) isLastItemToProcess(index, recordSize int) bool {
	return index == (recordSize - 1)
}

func (a *authorService) isAuthorNotAdded(authorsAddedMap map[string]bool, name string) bool {
	return !authorsAddedMap[name]
}
