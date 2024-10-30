BIN=doggo
ifeq ($(OS),Windows_NT)
	BIN=doggo.exe
endif

help:
ifeq ($(OS),Windows_NT)
	@echo "Windows help not supported"
else
	@echo "Usage:"
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ":" | sed -e 's/^/  /'
endif

## lint: Run the linter
lint:
	docker compose exec doggo.api ~/go/bin/golangci-lint run

## tests: Run unit tests
tests:
	go test -v ./...

## tests-coverage: Generate report for unit test coverage
tests-coverage:
	go test -coverprofile cover.out ./...
	go tool cover -html=cover.out

## tests-report: View report of all tests results
tests-report:
	go test ./... -json | tparse -all

## build: Builds the binary
build:
	go build -o bin/${BIN} main.go

## docker-modules: Cleans go.mod file and installs the required go build dependencies
docker-modules:
	go mod tidy

## compile: Compiles the production binaries
compile:
	echo "Compiling production binaries"
	GOOS=linux GOARCH=amd64 go build -o build/doggo-linux-amd64 main.go

all: build
