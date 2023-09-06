# Makefile

# Variables
APP_NAME := roadwarrior-backend
AWS_DEFAULT_PROFILE := localzone-project
GOPATH := $(shell go env GOPATH)
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=main

# Commands
all: clean build package

run-local:
	docker-compose -f ddb-local-docker-compose.yml up -d
	AWS_DEFAULT_PROFILE=$(AWS_DEFAULT_PROFILE) ENV=dev $(GOPATH)/bin/air

run-prod:
	docker-compose -f ddb-local-docker-compose.yml up -d
	AWS_DEFAULT_PROFILE=$(AWS_DEFAULT_PROFILE) ENV=prod $(GOPATH)/bin/air

build: 
	$(GOBUILD) ./...
	ENV=prod GOARCH=amd64 GOOS=linux $(GOBUILD) -o $(BINARY_NAME)

package:
	zip main.zip $(BINARY_NAME)
clean: 
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_NAME).zip