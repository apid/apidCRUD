#! /bin/bash

set -o physical
export HOME=~djfong
export GOROOT=/usr/local/go
export GOPATH=$HOME/edgex/go
PATH=:$GOROOT/bin:$GOPATH/bin:$PATH
