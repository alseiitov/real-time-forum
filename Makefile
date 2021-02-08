.PHONY: build
build:
	rm -rf build
	mkdir build
	go build -o api ./cmd/api/main.go
	go build -o client ./cmd/client/main.go

.PHONY: run-api
run-api:
ifdef config-path
	export CGO_ENABLED=1
	go run ./cmd/api/main.go -config-path="$(config-path)"
else
	export CGO_ENABLED=1
	go run ./cmd/api/main.go
endif
	 
.PHONY: run-client
run-client:
ifdef config-path
	go run ./cmd/client/main.go -config-path="$(config-path)"
else
	go run ./cmd/client/main.go 
endif
	 