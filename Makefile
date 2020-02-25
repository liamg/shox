default: build

build:
	./scripts/build.sh

test:
	go test -race ./...

run:
	go build ./cmd/shox
	./shox
