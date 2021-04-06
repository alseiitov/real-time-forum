include .env
ENV_VARS := JWT_SIGNING_KEY=${JWT_SIGNING_KEY} PASSWORD_SALT=${PASSWORD_SALT}

ifdef config-path
	ARGS := -config-path="$(config-path)"
endif

.PHONY:
.SILENT:

build:
	rm -rf build && mkdir build
	go build -o ./build/api ./cmd/api/main.go
	go build -o ./build/client ./cmd/client/main.go

run-api: build
	${ENV_VARS} ./build/api ${ARGS}

run-client: build
	${ENV_VARS} ./build/client ${ARGS}

lint:
	 golangci-lint run

swagger:
	swag init --parseDependency -g internal/app/app.go 