.PHONY: build
build:
	rm -rf build
	mkdir build
	go build -o build ./cmd/real-time-forum/main.go

.PHONY: run
run:
ifdef config-path
	export CGO_ENABLED=1
	go run ./cmd/real-time-forum/main.go -config-path="$(config-path)"
else
	export CGO_ENABLED=1
	go run ./cmd/real-time-forum/main.go
endif
	 
