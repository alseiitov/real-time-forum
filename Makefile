.PHONY: build run-api run-clent lint swagger
.SILENT:

build:
	rm -rf build
	go build -o ./build/api ./cmd/api/main.go
	go build -o ./build/client ./cmd/client/main.go

run-api: build
	./build/api

run-client: build
	./build/client

lint:
	golangci-lint run

swagger:
	swag init --parseDependency -g internal/app/app.go 