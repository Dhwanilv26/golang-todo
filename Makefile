run:
	go run ./cmd/api

dev:
	air

test:
	go test ./...

fmt:
	go fmt ./...

lint:
	golangci-lint run

migrate-up:
	migrate -path migrations -database "$(DATABASE_URL)" up

migrate-down:
	migrate -path migrations -database "$(DATABASE_URL)" down 1

create-migration:
	migrate create -ext sql -dir migrations -seq $(name)

psql:
	psql "$(DATABASE_URL)"