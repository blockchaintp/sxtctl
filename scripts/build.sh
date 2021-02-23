#!/bin/bash
set -xeuo pipefail
export DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

export TARGETDIR="$DIR/../target"
export CODE="$DIR/../cmd/sxtctl/main.go"
export CGO_ENABLED=0
export GO111MODULE=on
rm -rf $TARGETDIR
mkdir -p $TARGETDIR
for GOOS in darwin linux windows; do
  for GOARCH in 386 amd64; do
    export GOOS GOARCH
    echo "building sxtctl-$GOOS-$GOARCH"
    go build \
      -o "$TARGETDIR/sxtctl-$GOOS-$GOARCH" \
      -ldflags "-w -extldflags \"-static\"" \
      $CODE
  done
done
