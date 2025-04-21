.DEFAULT_GOAL := build

.PHONY: test fmt vet build
fmt:
	go fmt ./...

vet: fmt
	go vet ./...

test: vet
	go test -cover ./...

build: test
	go build -o flkcli main.go