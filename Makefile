.PHONY: build-cli build-server build-ui dev-server run-cli

build-cli:
	cd cli && go build -o ../bin/aiusage .

build-server:
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
