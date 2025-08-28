export PATH := $(GOPATH)/bin:$(PATH)
export GO111MODULE=on
LDFLAGS := -X 'main.buildTime=$(shell date -u +"%Y-%m-%d %H:%M:%S UTC")' -X 'main.gitCommit=$(shell git log --pretty=format:"%h" -1)' -X 'main.version=$(shell git describe --tags --exact-match 2>/dev/null || echo "DailyBuild")'
export CGO_LDFLAGS = "-lc++ -lduckdb_bundle -L/path/to/folder/with/lib"

all: fmt vet build

build: codeCli codeDebug codeAgent

release: build
	mkdir -p build/magicCodeAgent
	cp bin/magicCodeCli build/magicCodeAgent/magicCodeCli
	cp bin/magicCodeDebug build/magicCodeAgent/magicCodeDebug
	cp bin/magicCodeAgent build/magicCodeAgent/magicCodeAgent
	cd build && ls -l magicCodeAgent && zip -r magicCodeAgent.zip magicCodeAgent

install: fmt build
	cp bin/magicCodeCli /usr/local/bin/magicCodeCli
	cp bin/magicCodeDebug /usr/local/bin/magicCodeDebug
	cp bin/magicCodeAgent /usr/local/bin/magicCodeAgent

codeCli:
	env CGO_ENABLED=0 go build -ldflags "$(LDFLAGS)" -o bin/magicCodeCli ./application/magicCodeCli/cmd

codeDebug:
	env CGO_ENABLED=0 go build -ldflags "$(LDFLAGS)" -o bin/magicCodeDebug ./application/magicCodeDebug/cmd

codeAgent: downloadDuckdb
	go build -x -tags "no_duckdb_arrow duckdb_use_static_lib" -ldflags "$(LDFLAGS)" -o bin/magicCodeAgent ./application/magicCodeAgent/cmd

downloadDuckdb:
	./deps/download_duckdb.sh v1.3.2

fmt:
	go fmt ./...

vet:
	go vet ./...

clean:
	rm -rf ./bin & rm -rf build
