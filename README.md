## Golang-test
This app contains a CRUD to manage books and its authors.

### Cloud app APIs
- https://golang-test-books.herokuapp.com/swagger/index.html

### Requiriments
- [Golang 1.18+](https://go.dev/)
- [Postgres](https://www.postgresql.org/)
- [Docker](https://www.docker.com/)
- [Makefile](https://makefiletutorial.com/)

### Environment vars (for local environment):
```
      - DATABASE_URL=postgres://postgres:postgres@localhost:5432/postgres
      - PORT=3000
      - AUTHORS_FILE_PATH=./data/authorsreduced.csv
```

### How to use commands of Makefile

You can use Makefile to use commands of the app:

To run all services by docker
```
make services-all-up
```

To build app and run:
```
make run-services-dev
make build-run
```

To run for develop locally
```
make go-code
```

To run unit tests:
```
make tests
```

To run coverage tests:
```
make tests-coverage
```

To update docs from swagger

For generate docs of swagger, you need install [swagger-echo](https://github.com/swaggo/echo-swagger#start-using-it)

```
make swagger
```

### APIs
#### List all APIs
```
/swagger/index.html
```

### Tools and steps to built this app

- I used MacOS to develo app;
- I used VSCode IDE to code it;
- I used Postgres database;
- I used Echo for web framework;
- I used Swagger to generate docs of APIs;
- I used Gorm ORM;
- I used Docker;
- I used Makfile to improves experience on run app;

#### Some details about project building

- At the begin, I thought use go-routines/workers to import authors from csv file. But It was very slowly to import all authors one by one even using concurrency. So, I changed to use unique key name to avoid duplicate authors and I used insert in batch to import authors.
   - Observation: In branch 'turn-func-concurrence' I did some tests using concurrence/insert in batch to improve response time on import a big list of authors from csv file.

- On use docker-compose, the host to connect on database is different, when run only container of database, I need use localhost, but when is communicate between two containers, I need use postgres-go and create new network for this communication called golangtestdriver (You can see it on docker-compose.yml).

- To turn easier for manager relations of book/authors, I decided to use Gorm (I often use sqlx)

