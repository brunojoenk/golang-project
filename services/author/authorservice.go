package services

import (
	"encoding/csv"
	"github/brunojoenk/golang-test/models"
	authorrepo "github/brunojoenk/golang-test/repository/author"
	"os"
	"sync"

	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type GetAllAuthors func(filter models.GetAuthorsFilter) ([]models.Author, error)
type CreateAuthorInBatch func(author []*models.Author, batchSize int) error

type AuthorService struct {
	getAllAuthorsRepository GetAllAuthors
	createAuthorInBatchRepo CreateAuthorInBatch
}

// NewBookService Service Constructor
func NewAuthorService(db *gorm.DB) *AuthorService {
	repo := authorrepo.NewAuthorRepository(db)
	return &AuthorService{getAllAuthorsRepository: repo.GetAllAuthors, createAuthorInBatchRepo: repo.CreateAuthorInBatch}
}

func (a *AuthorService) GetAllAuthors(filter models.GetAuthorsFilter) (*models.AuthorResponseMetadata, error) {

	filter.Pagination.ValidValuesAndSetDefault()
	authors, err := a.getAllAuthorsRepository(filter)
	if err != nil {
		log.Error("Error on get all authors from repositoriy: ", err.Error())
		return nil, err
	}

	authorsResponse := make([]models.AuthorResponse, 0)
	for _, a := range authors {
		authorResponse := &models.AuthorResponse{
			Id:   a.Id,
			Name: a.Name,
		}
		authorsResponse = append(authorsResponse, *authorResponse)
	}

	authorResponseMetada := &models.AuthorResponseMetadata{
		Authors:    authorsResponse,
		Pagination: filter.Pagination,
	}

	return authorResponseMetada, nil
}

// Import all author using concurrence
func (a *AuthorService) ImportAuthorsFromCSVFile(file string) ([]string, error) {
	f, err := os.Open(file)

	if err != nil {
		log.Error("Error on open file: ", err.Error())
		return nil, err
	}

	defer f.Close()

	fcsv := csv.NewReader(f)
	fcsv.Comma = ';'
	authors := make([]*models.Author, 0)
	numWorkers := 20
	jobs := make(chan []*models.Author, numWorkers)
	res := make(chan []*models.Author)

	batchSize := 1000

	var wg sync.WaitGroup
	worker := func(jobs <-chan []*models.Author, results chan<- []*models.Author) error {
		for {
			select {
			case job, ok := <-jobs: // you must check for readable state of the channel.
				if !ok {
					return nil
				}
				err := a.createAuthorInBatchRepo(job, len(job))
				if err != nil {
					log.Error("Error on create author in batch repository: ", err.Error())
					return err
				}
				results <- job
			}
		}
	}

	var errOnBatch error
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
		mapper := make(map[string]bool, 0)
		rStr, err := fcsv.ReadAll()
		if err != nil {
			log.Error("Error on read all csv: ", err.Error())
			return
		}
		for _, record := range rStr {
			count := 0
			batch := make([]*models.Author, 0)
			for i, name := range record {
				if !mapper[name] {
					mapper[name] = true
					count++
					batch = append(batch, &models.Author{Name: name})
				}
				if count == batchSize || i == (len(record)-1) {
					jobs <- batch
					batch = make([]*models.Author, 0)
					count = 0
				}
			}
		}

		close(jobs) // close jobs to signal workers that no more job are incoming.
	}()

	go func() {
		wg.Wait()
		close(res) // when you close(res) it breaks the below loop.
	}()

	for r := range res {
		authors = append(authors, r...)
	}

	names := make([]string, 0)
	for _, a := range authors {
		names = append(names, a.Name)
	}

	return names, errOnBatch
}
