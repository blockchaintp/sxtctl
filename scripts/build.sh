#!/bin/bash
set -xeuo pipefail
export DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

export TARGETDIR="$DIR/../target"
export CODE="$DIR/../cmd/sxtctl/main.go"
export CGO_ENABLED=0
export GO111MODULE=on
rm -rf $TARGETDIR
mkdir -p $TARGETDIR
for ARCH in "darwin-amd64" "darwin-arm64" "linux-amd64" "windows-amd64"; do
    echo "building sxtctl-$ARCH"
    IFS='-'; arArch=($ARCH); unset IFS;
    export GOOS= ${arArch[0]}
    export GOARCH ${arArch[1]}

    go build \
      -o "$TARGETDIR/sxtctl-$ARCH" \
      -ldflags "-w -extldflags \"-static\"" \
      $CODE
done
