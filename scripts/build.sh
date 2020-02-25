#!/bin/bash
BINARY=shox
TAG=${TRAVIS_TAG:-development}
go mod download && go mod tidy
mkdir -p bin/darwin
GOOS=darwin GOARCH=amd64 go build -o bin/darwin/${BINARY}-darwin-amd64 ./cmd/shox/
mkdir -p bin/linux
GOOS=linux GOARCH=amd64 go build -o bin/linux/${BINARY}-linux-amd64 ./cmd/shox/
mkdir -p bin/windows
GOOS=windows GOARCH=amd64 go build -o bin/windows/${BINARY}-windows-amd64.exe ./cmd/shox/