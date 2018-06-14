# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_UNIX=$(BINARY_NAME)_unix

all: deps test client server
client:
		BINARY_NAME=client
		$(GOBUILD) -o $(BINARY_NAME) -v
server:
		BINARY_NAME=server
		$(GOBUILD) -o $(BINARY_NAME) -v
# build: 
		
test: 
		$(GOTEST) -v ./...
clean: 
		$(GOCLEAN)
		rm -f $(BINARY_NAME)
		rm -f $(BINARY_UNIX)
run:
		$(GOBUILD) -o $(BINARY_NAME) -v ./...
		./$(BINARY_NAME)
deps:
		# $(GOGET) github.com/

# Cross compilation
build-linux:
		CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_UNIX) -v