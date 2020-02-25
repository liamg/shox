default: build

build:
	go build ./cmd/shox

test:
	go test -race ./...

run: build
	./shox
