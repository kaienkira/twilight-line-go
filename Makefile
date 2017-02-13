PROJ_PATH = github.com/kaienkira/twilight-line-go

.PHONY: build-all build-server build-client

build-all: build-server build-client

build-server:
	go build -o bin/twilight-line-server $(PROJ_PATH)/server

build-client:
	go build -o bin/twilight-line-client $(PROJ_PATH)/client
