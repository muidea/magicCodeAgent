export PATH := $(GOPATH)/bin:$(PATH)
export GO111MODULE=on
LDFLAGS := -X 'main.buildTime=$(shell date -u --rfc-3339=seconds)' -X 'main.gitCommit=$(shell git log --pretty=format:"%h" -1)'


all: fmt vet build

build: cli codeAgent

release: build
	mkdir -p build/magicToolkit
	cp bin/magicCli build/magicToolkit/magicCli
	cp bin/magicCodeAgent build/magicToolkit/magicCodeAgent
	cd build && ls -l magicToolkit && zip -r magicToolkit.zip magicToolkit

install: build
	cp bin/magicCli /usr/local/bin/magicCli
	cp bin/magicCodeAgent /usr/local/bin/magicCodeAgent

cli:
	env CGO_ENABLED=0 go build -ldflags "$(LDFLAGS)" -o bin/magicCli ./application/magicCli/cmd

codeAgent:
	env CGO_ENABLED=0 go build -ldflags "$(LDFLAGS)" -o bin/magicCodeAgent ./application/magicCodeAgent/cmd


fmt:
	go fmt ./...

vet:
	go vet ./...

clean:
	rm -rf ./bin & rm -rf build
