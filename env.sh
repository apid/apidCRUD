#! /bin/bash

set -o physical
H=~djfong
# export GOROOT=${GOROOT:-/usr/local/go}
unset GOROOT
export GOPATH=${GOPATH:-$H/edgex/go}
PATH=:$H/homebrew/bin:$GOROOT/bin:$GOPATH/bin:$PATH
