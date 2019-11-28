# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=bakery
BINARY_UNIX=$(BINARY_NAME)_unix

all: test build
build: 
				$(GOBUILD) -o $(BINARY_NAME) -v
build_lambda:
				GOOS=linux $(GOBUILD) -o main -v
				zip function.zip main
test: 
				$(GOTEST) -v ./...
clean: 
				$(GOCLEAN)
				rm -f $(BINARY_NAME)
				rm -f $(BINARY_UNIX)
run:
				$(GOGET)
				$(GOBUILD) -o $(BINARY_NAME)
				./$(BINARY_NAME)

