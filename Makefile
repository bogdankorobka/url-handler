.PHONY: init run build swagger
.DEFAULT_GOAL := init

init: build run

build:
	go build -v ./cmd/handler

run:
	./handler

swagger:
	swag init -g internal/server/server.go
