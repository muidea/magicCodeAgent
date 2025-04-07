export PATH := $(GOPATH)/bin:$(PATH)
export GO111MODULE=on
LDFLAGS := -X 'main.buildTime=$(shell date -u --rfc-3339=seconds)' -X 'main.gitCommit=$(shell git log --pretty=format:"%h" -1)'


all: fmt vet build

build: toolkit

install: toolkit
	cp bin/magicToolkit /usr/local/bin/magicToolkit

toolkit:
	env CGO_ENABLED=0 go build -ldflags "$(LDFLAGS)" -o bin/magicToolkit ./application/magicToolkit/cmd

fmt:
	go fmt ./...

vet:
	go vet ./...

clean:
	rm -rf ./bin & rm -rf build
