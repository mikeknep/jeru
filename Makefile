.PHONY: test

build:
	go build -o out/jeru *.go

test:
	go test ./...
