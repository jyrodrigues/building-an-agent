
.DEFAULT_GOAL := run

.PHONY: generate build test

generate:
	baml-cli generate

build: generate
	go build *.go

run: generate
	go run *.go

# test: generate
# 	go test ./...

