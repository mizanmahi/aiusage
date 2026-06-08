.PHONY: build-cli build-server build-ui dev-server run-cli tidy test test-integration migrate migrate-down migrate-status

build-cli:
	mkdir -p bin
	cd cli && go build -o ../bin/aiusage .

build-server:
	mkdir -p bin
	cd server && go build -o ../bin/server .

build-ui:
	cd ui && pnpm run build

dev-server:
	cd server && set -a && { [ ! -f .env ] || . ./.env; } && set +a && go run .

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

test-integration:
	cd server && set -a && { [ ! -f .env ] || . ./.env; } && set +a && go test -tags=integration ./...

migrate:
	cd server && set -a && { [ ! -f .env ] || . ./.env; } && set +a && goose -dir migrations postgres "$$DATABASE_URL" up

migrate-down:
	cd server && set -a && { [ ! -f .env ] || . ./.env; } && set +a && goose -dir migrations postgres "$$DATABASE_URL" down

migrate-status:
	cd server && set -a && { [ ! -f .env ] || . ./.env; } && set +a && goose -dir migrations postgres "$$DATABASE_URL" status
