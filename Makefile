export AUTHORS_FILE_PATH=./data/authors.csv

services-all-up: 
	docker-compose up -d --build

services-down: 
	docker-compose down

run-main:
	go run main.go

run-services-dev:
	docker-compose -f docker-compose-dev.yml up

tests:
	go test  ./...

tests-coverage:
	go test -cover ./...

build:
	go build .

run-built:
	./golang-test

build-run:
	make build && make run-built

swagger:
	swag init

