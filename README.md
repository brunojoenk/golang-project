## Golang-test
This app contains a CRUD of books and APIs to manager authors of these books.

### Requiriments
- Golang 1.18+
- Docker (docker-compose)
- Swagger (https://github.com/swaggo/echo-swagger)
- Makefile

### Environment vars (for local environment):
```
      - POSTGRES_HOST=lcoalhost
      - POSTGRES_PORT=5432
      - POSTGRES_USER=postgres
      - POSTGRES_PASSWORD=postgres
      - POSTGRES_DBNAME=postgres
      - SERVER_PORT=3000
```

### How to run locally

You can use Makefile to use commands of the app:

To run all services by docker
```
make services-all-up
```

To build and run project:
```
make run-services-dev
make build-run
```

To run main.go:
```
make run-services-dev
make run-main
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
```
make swagger
```

### APIs
#### List of APIs
```
/swagger/index.html
```

### Tools and steps to built this app

####

- Used MacOS to develo app
- Used Echo for web framework;
- Used Swagger to generate and able APIs;
- Used Gorm ORM
- Used Docker
- Used Make

#### Difficults

- At the begin, I thought use go-routines/workers to import authors from csv file. But It was slowly to import all authors. So, I changed to use unique key name to avoid duplicate authors and I used insert in batch to import authors.

- On use docker-compose, the host to connect on database is different, when run only container of database, I need use localhost, but when is communicate between two containers, I need use postgres-go and create new network for this communication called golangtestdriver (You can see it on docker-compose.yml).

- To turn easier for manager relations of book/authors, I decided to use Gorm.