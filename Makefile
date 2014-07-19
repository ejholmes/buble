.PHONY: test example

test:
	go test ./...

example: example/main.go
	go build -o example/example ./example
