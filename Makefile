.PHONY: run build test test-cover lint mocks docker-build docker-run tidy fmt

## Local development ------------------------------------------------------

run: ## Run the API locally (reads ./.env, see .env.example)
	go run .

build: ## Compile the binary
	go build -o epa-standards-backend .

fmt: ## Format the code
	gofmt -w .

tidy: ## Sync go.mod/go.sum with imports
	go mod tidy

## Testing & linting -------------------------------------------------------

test: ## Run all unit tests
	go test ./... -v

test-cover: ## Run tests and print coverage per package
	go test ./... -cover

test-cover-html: ## Run tests and open an HTML coverage report
	go test ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html

lint: ## Run golangci-lint (install: https://golangci-lint.run/usage/install/)
	golangci-lint run

## Mocks --------------------------------------------------------------------

mocks: ## Regenerate every mock in mocks/ from internal/pkg/ports (see .mockery.yaml)
	mockery

## Docker ---------------------------------------------------------------

docker-build: ## Build the production image
	docker build -t epa-standards-backend:latest .

docker-run: ## Run the image locally with .env
	docker run --env-file .env -p 8080:8080 epa-standards-backend:latest
