run:
	go run ./cmd/api

test:
	go test ./...

lint:
	golangci-lint run

fmt:
	go fmt ./...

dev:
	air

migrate-create:
	migrate create -ext sql -dir migrations -rec 

migrate-up:
	migrate -path migrations -database $$DATABASE_URL up

migrate-down:
	migrate -path migrations -database $$DATABASE_URL down

migrate-force:
	migrate -path migrations -database $$DATABASE_URL force 1
