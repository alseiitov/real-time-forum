.PHONY: build
build:
	rm -rf build
	mkdir build
	go build -o build ./server/main.go

.PHONY: run
run:
ifdef config-path
	go run ./server/main.go -config-path="$(config-path)"
else
	go run ./server/main.go
endif
	 
