.PHONY: build
build:
	go build -v ./cmd/apiserver

.PHONY: test
test:
	go test -v -race -timeout 30s ./...

.PHONY: db up
migrate up:
	migrate -path db/migration -database "postgresql://oleg:password@localhost:5432/restapi_dev?sslmode=disable" -verbose up

.PHONY: db down
migrate down:
	migrate -path db/migration -database "postgresql://oleg:password@localhost:5432/restapi_dev?sslmode=disable" -verbose down

.DEFAULT_GOAL := build
