#!/bin/bash

proj_path=.
os_list="linux windows darwin freebsd"
arch_list="amd64 386"

for os in $os_list
do
    if [ "$os" = "windows" ]
    then
        ext=.exe
    else
        ext=
    fi

    for arch in $arch_list
    do
        GOOS=$os GOARCH=$arch go build \
            -o bin/twilight-line-go-server-$os-$arch$ext $proj_path/server
        GOOS=$os GOARCH=$arch go build \
            -o bin/twilight-line-go-client-$os-$arch$ext $proj_path/client
    done
done
