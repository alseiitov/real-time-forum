.PHONY: build
build:
	rm -rf build/*
	go build -o build ./cmd/real-time-forum/main.go

.PHONY: run
run:
ifdef config-path
	go run ./cmd/real-time-forum/main.go -config-path="$(config-path)"
else
	go run ./cmd/real-time-forum/main.go
endif
	 