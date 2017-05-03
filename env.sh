#! /bin/bash

set -o physical
export GOROOT=${GOROOT:-/usr/local/go}
export GOPATH=${GOPATH:-$HOME/edgex/go}
PATH=:$GOROOT/bin:$GOPATH/bin:$PATH
