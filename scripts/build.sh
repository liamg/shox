#!/bin/bash
BINARY=shox
TAG=${TRAVIS_TAG:-development}
mkdir -p bin/darwin
GOOS=darwin GOARCH=amd64 go build -mod=vendor -o bin/darwin/${BINARY}-darwin-amd64 ./cmd/shox/
mkdir -p bin/linux
GOOS=linux GOARCH=amd64 go build -mod=vendor -o bin/linux/${BINARY}-linux-amd64 ./cmd/shox/
# compilation over ARM architectures
GOOS=linux GOARCH=arm go build -mod=vendor -o bin/linux/${BINARY}-linux-arm ./cmd/shox/
