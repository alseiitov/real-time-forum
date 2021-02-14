include .env

.PHONY: build
build:
	@rm -rf build
	@mkdir build
	go build -o ./build/api ./cmd/api/main.go
	go build -o ./build/client ./cmd/client/main.go

.PHONY: run-api
run-api:
ifdef config-path
	@CGO_ENABLED=1 \
	SECRET_KEY=${SECRET_KEY} \
	go run ./cmd/api/main.go -config-path="${config-path}" 
else
	@CGO_ENABLED=1 \
	SECRET_KEY=${SECRET_KEY} \
	go run ./cmd/api/main.go
endif
	 
.PHONY: run-client
run-client:
ifdef config-path
	go run ./cmd/client/main.go -config-path="$(config-path)"
else
	go run ./cmd/client/main.go 
endif
	 