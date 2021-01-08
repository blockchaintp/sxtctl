#!/bin/bash
set -xeuo pipefail
export DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

mkdir -p "$DIR/../target/{Linux-x86_64,Linux-arm64,Darwin-x86_64}"
GOOS=linux GOARCH=amd64 GO111MODULE=on go build -o "$DIR/../target/Linux-x86_64/sxtctl" "$DIR/../cmd/sxtctl/main.go"
#GOOS=linux GOARCH=arm64 GO111MODULE=on go build -o "$DIR/../target/Linux-arm64/sxtctl" "$DIR/../cmd/sxtctl/main.go"
#GOOS=darwin GOARCH=amd64 GO111MODULE=on go build -o "$DIR/../target/Darwin-x86_64/sxtctl" "$DIR/../cmd/sxtctl/main.go"

