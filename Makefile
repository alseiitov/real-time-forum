.PHONY:
include .env

build:
	@rm -rf build
	@mkdir build
	@go build -o ./build/api ./cmd/api/main.go
	@go build -o ./build/client ./cmd/client/main.go

run-api:
ifdef config-path
	@CGO_ENABLED=1 \
	JWT_SIGNING_KEY=${JWT_SIGNING_KEY} \
	PASSWORD_SALT=${PASSWORD_SALT} \
	go run ./cmd/api/main.go -config-path="${config-path}" 
else
	@CGO_ENABLED=1 \
	JWT_SIGNING_KEY=${JWT_SIGNING_KEY} \
	PASSWORD_SALT=${PASSWORD_SALT} \
	go run ./cmd/api/main.go
endif

run-client:
ifdef config-path
	@go run ./cmd/client/main.go -config-path="$(config-path)"
else
	@go run ./cmd/client/main.go 
endif
	 
lint:
	 golangci-lint run

swagger:
	swag init --parseDependency -g internal/app/app.go 