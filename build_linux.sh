#!/bin/bash

proj_path=.
os=linux
arch=amd64

GOOS="$os" GOARCH="$arch" go build \
    -o bin/twilight-line-go-server-"$os"-"$arch" "$proj_path"/server
if [ $? -ne 0 ]; then exit 1; fi

GOOS="$os" GOARCH="$arch" go build \
    -o bin/twilight-line-go-client-"$os"-"$arch" "$proj_path"/client
if [ $? -ne 0 ]; then exit 1; fi

exit 0
