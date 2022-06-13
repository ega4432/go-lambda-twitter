.PHONY: build
buid: fmt
	go build -o main

.PHONY: fmt
fmt:
	go fmt ./...

.PHONY: install
install:
	go get -u ./...

.PHONY: test
test:
	go test -v ./...
