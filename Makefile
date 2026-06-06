.PHONY: build-cli build-server build-ui dev-server run-cli tidy test migrate

build-cli:
	mkdir -p bin
	cd cli && go build -o ../bin/aiusage .

build-server:
	mkdir -p bin
	cd server && go build -o ../bin/server .

build-ui:
	cd ui && pnpm run build

dev-server:
	cd server && go run .

run-cli:
	cd cli && go run . $(ARGS)

tidy:
	cd types && go mod tidy
	cd cli && go mod tidy
	cd server && go mod tidy

test:
	cd types && go test ./...
	cd cli && go test ./...
	cd server && go test ./...

migrate:
	psql $$DATABASE_URL -f server/migrations/001_initial.sql
