default: build

build:
	./scripts/build.sh

test:
	go test -mod=vendor -race ./...

run:
	go build -mod=vendor ./cmd/shox
	./shox
