export PATH := $(GOPATH)/bin:$(PATH)
export GO111MODULE=on
LDFLAGS := -X 'main.buildTime=$(shell date -u --rfc-3339=seconds)' -X 'main.gitCommit=$(shell git log --pretty=format:"%h" -1)'


all: fmt vet build

build: cli codeAgent

release: build
	mkdir -p build/magicCodeAgent
	cp bin/magicCodeCli build/magicCodeAgent/magicCodeCli
	cp bin/magicCodeAgent build/magicCodeAgent/magicCodeAgent
	cd build && ls -l magicCodeAgent && zip -r magicCodeAgent.zip magicCodeAgent

install: build
	cp bin/magicCodeCli /usr/local/bin/magicCodeCli
	cp bin/magicCodeAgent /usr/local/bin/magicCodeAgent

cli:
	env CGO_ENABLED=0 go build -ldflags "$(LDFLAGS)" -o bin/magicCodeCli ./application/magicCodeCli/cmd

codeAgent:
	env CGO_ENABLED=0 go build -ldflags "$(LDFLAGS)" -o bin/magicCodeAgent ./application/magicCodeAgent/cmd


fmt:
	go fmt ./...

vet:
	go vet ./...

clean:
	rm -rf ./bin & rm -rf build
