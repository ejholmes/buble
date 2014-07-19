.PHONY: test example

test:
	go test -race ./...

example: example/main.go
	go build -o example/example ./example
