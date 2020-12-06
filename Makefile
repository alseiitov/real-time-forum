.PHONY: build
build:
	rm -rf build/*
	go build -o build cmd/real-time-forum/main.go

.PHONY: build
run:
	go run cmd/real-time-forum/main.go