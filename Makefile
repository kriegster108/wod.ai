ifneq (,$(wildcard ./.env))
    include .env
    export
endif

lint: 
	go vet

build:
	go build -o bin/main main.go

run:
	go run main.go