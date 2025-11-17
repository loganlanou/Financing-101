BINARY:=logans3d
DB_PATH:=data/app.db
DB_DSN:=sqlite://$(DB_PATH)

.PHONY: all build run clean generate templ sqlc migrate-up migrate-down migrate-status test css fmt lint seed deps

all: build

deps:
	go mod tidy
	npm install

build: generate css
	GOOS=$$(go env GOOS) GOARCH=$$(go env GOARCH) go build -o $(BINARY) ./cmd/web

run: build
	./$(BINARY)

clean:
	rm -f $(BINARY)

fmt:
	go fmt ./...

lint:
	go vet ./...

css:
	npm run build:css

seed:
	go run ./cmd/tools/seed

test:
	go test ./...

migrate-up:
	goose -dir db/migrations sqlite3 $(DB_PATH) up

migrate-down:
	goose -dir db/migrations sqlite3 $(DB_PATH) down

migrate-status:
	goose -dir db/migrations sqlite3 $(DB_PATH) status

sqlc:
	sqlc generate

templ:
	templ generate

generate: templ sqlc
