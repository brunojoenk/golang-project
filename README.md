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

Run all services by docker
```
make services-all-up
```

Run for develop locally
```
make go-code
```

Build app and run:
```
make build-run
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

## Tools and Steps Used in Building This App

- **Operating System:** macOS was utilized for development.
- **IDE:** VSCode served as the Integrated Development Environment for coding.
- **Database:** PostgreSQL was chosen as the database management system.
- **Web Framework:** Echo was the web framework employed.
- **API Documentation:** Swagger was utilized to generate documentation for the APIs.
- **ORM:** Gorm was used as the Object-Relational Mapping tool.
- **Containerization:** Docker facilitated containerization.
- **Makefile:** A Makefile was implemented to enhance the experience of running the app.

## Project Building Details

Initially, I contemplated employing Go routines/workers to import authors from a CSV file. However, even with concurrency, the process proved sluggish. Consequently, I opted to utilize a unique key name to prevent duplicate authors and implemented batch insertion to expedite the import process.

**Observation:** In the 'turn-func-concurrence' branch, I conducted tests leveraging concurrency and batch insertion to enhance the response time when importing a large list of authors from a CSV file.

When utilizing Docker Compose, the host for connecting to the database differs. When running only the database container, 'localhost' is used. However, when communication occurs between two containers, 'postgres-go' is employed. A new network named 'golangtestdriver' was created for this communication, as specified in the `docker-compose.yml` file.

To simplify the management of relationships between books and authors, Gorm was chosen over my usual choice, sqlx.


