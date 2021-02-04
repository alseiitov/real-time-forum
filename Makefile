.PHONY: build
build:
	rm -rf build
	mkdir build
	go build -o api ./cmd/app/main.go
	go build -o client ./cmd/web/web-server/main.go

.PHONY: run-backend
run-api:
ifdef config-path
	export CGO_ENABLED=1
	go run ./cmd/api/main.go -config-path="$(config-path)"
else
	export CGO_ENABLED=1
	go run ./cmd/api/main.go
endif
	 
.PHONY: run-frontend
run-client:
ifdef config-path
	go run ./cmd/client/main.go -config-path="$(config-path)"
else
	go run ./cmd/client/main.go 
endif
	 