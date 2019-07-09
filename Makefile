PROJ_PATH = .

.PHONY: build-all build-server build-client

build-all: build-server build-client

build-server:
	go build -o bin/twilight-line-server "$(PROJ_PATH)"/server

build-client:
	go build -o bin/twilight-line-client "$(PROJ_PATH)"/client

run-server:
	bin/twilight-line-server -e etc/server_config.json

run-client:
	bin/twilight-line-client -e etc/client_config.json

check-format:
	@find . -name '*.go' -exec gofmt -d {} \; -print

do-format:
	@find . -name '*.go' -exec gofmt -w {} \; -print

clean:
	@rm -f build/*

bs: build-server
bc: build-client
rs: run-server
rc: run-client
fmt: check-format
dofmt: do-format
