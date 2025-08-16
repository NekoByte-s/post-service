# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod

# Binary names
BINARY_NAME=modular
BINARY_UNIX=$(BINARY_NAME)_unix

.PHONY: all build clean test deps swagger docker-build

all: test build

build:
	$(GOBUILD) -o $(BINARY_NAME) -v ./cmd/$(BINARY_NAME)

build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_UNIX) -v ./cmd/$(BINARY_NAME)

clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_UNIX)

test:
	$(GOTEST) -v ./...

deps:
	$(GOMOD) download
	$(GOMOD) verify

tidy:
	$(GOMOD) tidy

swagger:
	swag init -g cmd/$(BINARY_NAME)/main.go -o docs

run:
	$(GOBUILD) -o $(BINARY_NAME) -v ./cmd/$(BINARY_NAME)
	./$(BINARY_NAME)

docker-build:
	docker build -t post-service:latest .

docker-run:
	docker run -p 8080:8080 post-service:latest