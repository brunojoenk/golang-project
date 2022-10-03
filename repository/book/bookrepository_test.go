package repository

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"github/brunojoenk/golang-test/models/dtos"
	"github/brunojoenk/golang-test/models/entities"
	"regexp"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-test/deep"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Suite struct {
	suite.Suite
	DB   *gorm.DB
	mock sqlmock.Sqlmock

	repository *BookRepository
}

func (s *Suite) AfterTest(_, _ string) {
	require.NoError(s.T(), s.mock.ExpectationsWereMet())
}

func TestInit(t *testing.T) {
	suite.Run(t, new(Suite))
}

func (s *Suite) SetupSuite() {
	var (
		db  *sql.DB
		err error
	)

	db, s.mock, err = sqlmock.New()
	require.NoError(s.T(), err)

	dialector := postgres.New(postgres.Config{
		DSN:                  "sqlmock_db_0",
		DriverName:           "postgres",
		Conn:                 db,
		PreferSimpleProtocol: true,
	})
	s.DB, err = gorm.Open(dialector, &gorm.Config{})
	require.NoError(s.T(), err)

	s.repository = &BookRepository{s.DB}
}

func (s *Suite) Test_repository_Create_Book() {
	var (
		bookId          = 1
		name            = "test-name"
		edition         = "edition"
		publicationYear = 2022

		authorId   = 2
		authorName = "brad"
	)

	s.mock.ExpectBegin()

	s.mock.ExpectQuery(regexp.QuoteMeta(
		`INSERT INTO "books" ("name","edition","publication_year") VALUES ($1,$2,$3) RETURNING "id"`)).
		WithArgs(name, edition, publicationYear).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).
			AddRow(bookId))

	s.mock.ExpectQuery(regexp.QuoteMeta(
		`INSERT INTO "authors" ("name","id") VALUES ($1,$2) ON CONFLICT DO NOTHING RETURNING "id"`)).
		WithArgs(authorName, authorId).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).
			AddRow(authorId))

	s.mock.ExpectExec(regexp.QuoteMeta(
		`INSERT INTO "author_book" ("book_id","author_id") VALUES ($1,$2) ON CONFLICT DO NOTHING`)).
		WithArgs(bookId, authorId).WillReturnResult(driver.ResultNoRows)

	s.mock.ExpectCommit()

	err := s.repository.CreateBook(&entities.Book{
		Name:            name,
		Edition:         edition,
		PublicationYear: publicationYear,
		Authors:         []*entities.Author{{Id: authorId, Name: authorName}}})

	require.NoError(s.T(), err)
}

func (s *Suite) Test_repository_Update_Book() {
	var (
		bookId          = 1
		name            = "test-name"
		edition         = "edition"
		publicationYear = 2022

		authorId   = 2
		authorName = "brad"
	)

	//s.mock.MatchExpectationsInOrder(false)
	s.mock.ExpectBegin()

	s.mock.ExpectExec(regexp.QuoteMeta(
		`DELETE FROM "author_book" WHERE "author_book"."book_id" = $1`)).
		WithArgs(bookId).WillReturnResult(sqlmock.NewResult(int64(bookId), 1))

	s.mock.ExpectCommit()

	s.mock.ExpectBegin()

	s.mock.ExpectExec(regexp.QuoteMeta(
		`UPDATE "books" SET "name"=$1,"edition"=$2,"publication_year"=$3 WHERE "id" = $4`)).
		WithArgs(name, edition, publicationYear, bookId).WillReturnResult(sqlmock.NewResult(int64(bookId), 1))

	s.mock.ExpectQuery(regexp.QuoteMeta(
		`INSERT INTO "authors" ("name","id") VALUES ($1,$2) ON CONFLICT DO NOTHING RETURNING "id"`)).
		WithArgs(authorName, authorId).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).
			AddRow(authorId))

	s.mock.ExpectExec(regexp.QuoteMeta(
		`INSERT INTO "author_book" ("book_id","author_id") VALUES ($1,$2) ON CONFLICT DO NOTHING`)).
		WithArgs(bookId, authorId).WillReturnResult(sqlmock.NewResult(int64(bookId), int64(authorId)))

	s.mock.ExpectCommit()

	err := s.repository.UpdateBook(&entities.Book{
		Id:              bookId,
		Name:            name,
		Edition:         edition,
		PublicationYear: publicationYear}, []*entities.Author{{Id: authorId, Name: authorName}})

	require.NoError(s.T(), err)
}

func (s *Suite) Test_repository_Update_Book_Error_On_Clear() {
	var (
		bookId          = 1
		name            = "test-name"
		edition         = "edition"
		publicationYear = 2022

		authorId   = 2
		authorName = "brad"
	)

	//s.mock.MatchExpectationsInOrder(false)
	s.mock.ExpectBegin()

	s.mock.ExpectExec(regexp.QuoteMeta(
		`DELETE FROM "author_book" WHERE "author_book"."book_id" = $1`)).
		WithArgs(bookId).WillReturnError(context.Canceled)

	s.mock.ExpectRollback()

	err := s.repository.UpdateBook(&entities.Book{
		Id:              bookId,
		Name:            name,
		Edition:         edition,
		PublicationYear: publicationYear}, []*entities.Author{{Id: authorId, Name: authorName}})

	require.Error(s.T(), err)
}

func (s *Suite) Test_repository_Update_Book_Error_On_Save() {
	var (
		bookId          = 1
		name            = "test-name"
		edition         = "edition"
		publicationYear = 2022

		authorId   = 2
		authorName = "brad"
	)

	//s.mock.MatchExpectationsInOrder(false)
	s.mock.ExpectBegin()

	s.mock.ExpectExec(regexp.QuoteMeta(
		`DELETE FROM "author_book" WHERE "author_book"."book_id" = $1`)).
		WithArgs(bookId).WillReturnResult(sqlmock.NewResult(int64(bookId), 1))

	s.mock.ExpectCommit()

	s.mock.ExpectBegin()

	s.mock.ExpectExec(regexp.QuoteMeta(
		`UPDATE "books" SET "name"=$1,"edition"=$2,"publication_year"=$3 WHERE "id" = $4`)).
		WithArgs(name, edition, publicationYear, bookId).WillReturnError(context.Canceled)

	s.mock.ExpectRollback()

	err := s.repository.UpdateBook(&entities.Book{
		Id:              bookId,
		Name:            name,
		Edition:         edition,
		PublicationYear: publicationYear}, []*entities.Author{{Id: authorId, Name: authorName}})

	require.Error(s.T(), err)
}

func (s *Suite) Test_repository_Create_Book_Error() {
	var (
		name            = "test-name"
		edition         = "edition"
		publicationYear = 2022
	)

	s.mock.ExpectBegin()

	s.mock.ExpectQuery(regexp.QuoteMeta(
		`INSERT INTO "books" ("name","edition","publication_year") VALUES ($1,$2,$3) RETURNING "id"`)).
		WithArgs(name, edition, publicationYear).
		WillReturnError(context.Canceled)

	s.mock.ExpectRollback()

	err := s.repository.CreateBook(&entities.Book{
		Name:            name,
		Edition:         edition,
		PublicationYear: publicationYear})

	require.Error(s.T(), err)
}

func (s *Suite) Test_repository_Get_Book() {
	var (
		id         = 1
		name       = "test-name"
		authorId   = 2
		authorName = "author-name"
	)

	s.mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "books" WHERE "books"."id" = $1 ORDER BY "books"."id" LIMIT 1`)).
		WithArgs(id).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).
			AddRow(id, name))

	s.mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "author_book" WHERE "author_book"."book_id" = $1`)).
		WithArgs(id).
		WillReturnRows(sqlmock.NewRows([]string{"book_id", "author_id"}).
			AddRow(id, authorId))

	s.mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "authors" WHERE "authors"."id" = $1`)).
		WithArgs(authorId).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).
			AddRow(authorId, authorName))

	res, err := s.repository.GetBook(id)

	require.NoError(s.T(), err)
	require.Nil(s.T(), deep.Equal(&entities.Book{
		Id:   id,
		Name: name,
		Authors: []*entities.Author{
			{Id: authorId, Name: authorName}}}, res))
}

func (s *Suite) Test_repository_Get_Book_Error() {
	var (
		id = 1
	)

	s.mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "books" WHERE "books"."id" = $1 ORDER BY "books"."id" LIMIT 1`)).
		WithArgs(id).
		WillReturnError(context.Canceled)

	_, err := s.repository.GetBook(id)

	require.Error(s.T(), err)
}

func (s *Suite) Test_repository_Get_All_Books() {
	var (
		id              = 1
		name            = "test-name"
		edition         = "firtEsdiction"
		publicationYear = 2022
		authorId        = 2
		authorName      = "author-name"
	)

	s.mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT "books"."id","books"."name","books"."edition","books"."publication_year" 
		FROM "books" 
		JOIN author_book ON author_book.book_id = books.id 
		JOIN authors ON authors.id = author_book.author_id 
		WHERE LOWER(authors.name) LIKE $1 
		AND LOWER(books.name) LIKE $2
		AND LOWER(books.edition) LIKE $3 
		AND books.publication_year = $4`)).
		WithArgs("%"+strings.ToLower(authorName)+"%", "%"+strings.ToLower(name)+"%", "%"+strings.ToLower(edition)+"%", publicationYear).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).
			AddRow(id, name))

	s.mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "author_book" WHERE "author_book"."book_id" = $1`)).
		WithArgs(id).
		WillReturnRows(sqlmock.NewRows([]string{"book_id", "author_id"}).
			AddRow(id, authorId))

	s.mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "authors" WHERE "authors"."id" = $1`)).
		WithArgs(authorId).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).
			AddRow(authorId, authorName))

	res, err := s.repository.GetAllBooks(dtos.GetBooksFilter{Name: name, Edition: edition, PublicationYear: publicationYear, Author: authorName})

	require.NoError(s.T(), err)
	require.Nil(s.T(), deep.Equal([]entities.Book{{
		Id:   id,
		Name: name,
		Authors: []*entities.Author{
			{Id: authorId, Name: authorName}}}}, res))
}

func (s *Suite) Test_repository_Get_All_Books_Error() {
	var (
		name            = "test-name"
		edition         = "firtEsdiction"
		publicationYear = 2022
		authorName      = "author-name"
	)

	s.mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT "books"."id","books"."name","books"."edition","books"."publication_year" 
		FROM "books" 
		JOIN author_book ON author_book.book_id = books.id 
		JOIN authors ON authors.id = author_book.author_id 
		WHERE LOWER(authors.name) LIKE $1 
		AND LOWER(books.name) LIKE $2
		AND LOWER(books.edition) LIKE $3 
		AND books.publication_year = $4`)).
		WithArgs("%"+strings.ToLower(authorName)+"%", "%"+strings.ToLower(name)+"%", "%"+strings.ToLower(edition)+"%", publicationYear).
		WillReturnError(context.Canceled)

	_, err := s.repository.GetAllBooks(dtos.GetBooksFilter{Name: name, Edition: edition, PublicationYear: publicationYear, Author: authorName})

	require.Error(s.T(), err)

}

func (s *Suite) Test_repository_Delete_Book() {
	var (
		bookId   = 1
		bookName = "book"
	)

	s.mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "books" WHERE "books"."id" = $1 ORDER BY "books"."id" LIMIT 1`)).
		WithArgs(bookId).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).
			AddRow(bookId, bookName))

	s.mock.ExpectExec(regexp.QuoteMeta(
		`DELETE FROM author_book WHERE author_book.book_id = $1`)).
		WithArgs(bookId).WillReturnResult(sqlmock.NewResult(int64(bookId), 1))

	s.mock.ExpectBegin()

	s.mock.ExpectExec(regexp.QuoteMeta(
		`DELETE FROM "books" WHERE "books"."id" = $1`)).
		WithArgs(bookId).WillReturnResult(sqlmock.NewResult(int64(bookId), 1))

	s.mock.ExpectCommit()

	err := s.repository.DeleteBook(bookId)

	require.NoError(s.T(), err)
}

func (s *Suite) Test_repository_Delete_Book_Error_On_Select() {
	var (
		bookId = 1
	)

	s.mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "books" WHERE "books"."id" = $1 ORDER BY "books"."id" LIMIT 1`)).
		WithArgs(bookId).
		WillReturnError(context.Canceled)

	err := s.repository.DeleteBook(bookId)

	require.Error(s.T(), err)
}

func (s *Suite) Test_repository_Delete_Book_Error_On_Delete_Foreign_key() {
	var (
		bookId   = 1
		bookName = "book"
	)

	s.mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "books" WHERE "books"."id" = $1 ORDER BY "books"."id" LIMIT 1`)).
		WithArgs(bookId).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).
			AddRow(bookId, bookName))

	s.mock.ExpectExec(regexp.QuoteMeta(
		`DELETE FROM author_book WHERE author_book.book_id = $1`)).
		WithArgs(bookId).
		WillReturnError(context.Canceled)

	err := s.repository.DeleteBook(bookId)

	require.Error(s.T(), err)
}

func (s *Suite) Test_repository_Delete_Book_Error_On_Delete_Book() {
	var (
		bookId   = 1
		bookName = "book"
	)

	s.mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "books" WHERE "books"."id" = $1 ORDER BY "books"."id" LIMIT 1`)).
		WithArgs(bookId).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).
			AddRow(bookId, bookName))

	s.mock.ExpectExec(regexp.QuoteMeta(
		`DELETE FROM author_book WHERE author_book.book_id = $1`)).
		WithArgs(bookId).WillReturnResult(sqlmock.NewResult(int64(bookId), 1))

	s.mock.ExpectBegin()

	s.mock.ExpectExec(regexp.QuoteMeta(
		`DELETE FROM "books" WHERE "books"."id" = $1`)).
		WithArgs(bookId).WillReturnError(context.Canceled)

	s.mock.ExpectRollback()

	err := s.repository.DeleteBook(bookId)

	require.Error(s.T(), err)
}
