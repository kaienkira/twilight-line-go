#!/bin/bash

os=darwin
arch=arm64

cd client && \
GOOS="$os" GOARCH="$arch" go build \
    -o ../bin/twilight-line-go-client-"$os"-"$arch" && \
cd - >/dev/null
if [ $? -ne 0 ]; then exit 1; fi

cd server && \
GOOS="$os" GOARCH="$arch" go build \
    -o ../bin/twilight-line-go-server-"$os"-"$arch" && \
cd - >/dev/null
if [ $? -ne 0 ]; then exit 1; fi

exit 0
