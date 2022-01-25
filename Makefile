PROJECTNAME=$(shell basename "$(PWD)")

GOBASE=$(shell pwd)
GOBIN=$(GOBASE)/bin

all: build

install:
	go mod download

build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o $(GOBIN)/app ./main.go || exit

dev:
	go build -o $(GOBIN)/app ./main.go || exit
	$(GOBIN)/app