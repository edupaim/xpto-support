APPLICATION_NAME := xpto-support
BIN_NAME=${APPLICATION_NAME}

.PHONY: help

help: ## This help.
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

.DEFAULT_GOAL := help

default: help

init-stack:
	docker-compose up

build: ## Build project for development
	@echo "building ${BIN_NAME}"
	GOOS=linux GOARCH=amd64 go build -o bin/${BIN_NAME} ./

get-deps: ## Install projects dependencies with Go Module
	go mod tidy

docker-build: build  ## Build docker image
	docker build -t ${APPLICATION_NAME}:0.0.0 ./

run-test:  ## Run project tests
	mkdir -p ./test/cover
	go test ./... -race -coverpkg=./... -coverprofile=./test/cover/cover.out
	go tool cover -html=./test/cover/cover.out -o ./test/cover/cover.html
