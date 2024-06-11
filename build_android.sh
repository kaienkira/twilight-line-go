#!/bin/bash

os=android
android_compiler_dir=\
$ANDROID_HOME/ndk/26.3.11579264/toolchains/llvm/prebuilt/linux-x86_64/bin

build()
{
    local android_arch=$1
    local go_arch=$2
    local android_sdk_version=$3
    local compiler=$(printf '%s/%s-linux-%s-clang' \
          "$android_compiler_dir" \
          "$android_arch" "$android_sdk_version")

    cd client && \
    CC="$compiler" GOOS="$os" GOARCH="$go_arch" CGO_ENABLED=1 go build \
        -o ../bin/twilight-line-go-client-"$os"-"$android_arch" && \
    cd - >/dev/null
    if [ $? -ne 0 ]; then return 1; fi

    cd server && \
    CC="$compiler" GOOS="$os" GOARCH="$go_arch" CGO_ENABLED=1 go build \
        -o ../bin/twilight-line-go-server-"$os"-"$android_arch" && \
    cd - >/dev/null
    if [ $? -ne 0 ]; then return 1; fi

    return 0
}

build aarch64 arm64 android34
if [ $? -ne 0 ]; then exit 1; fi
build armv7a arm androideabi34
if [ $? -ne 0 ]; then exit 1; fi
build x86_64 amd64 android34
if [ $? -ne 0 ]; then exit 1; fi
build i686 386 android34
if [ $? -ne 0 ]; then exit 1; fi

exit 0
