#!/usr/bin/env bash

modules()
{
	go list ./... | grep -v vendor
}


# ----- start of mainline code

cdir=cov
ctxt=$cdir/coverage.txt

mkdir -p "$cdir"

set -e
echo "mode: atomic" > "$ctxt"

/bin/rm -f "$prof"

for d in $(modules); do
    prof=$cdir/profile.out
    go test "-coverprofile=$prof" -covermode=atomic $d
    if [ -f "$prof" ]; then
        head -2 "$prof" >> "$ctxt"
        rm "$prof"
    fi
done
go tool cover "-html=$ctxt" -o "$cdir/cover.html"
