#! /bin/bash

ID=${1:-25}
TABLE=bundles
out=$(./appcurl.sh DELETE "db/_table/$TABLE/$ID")
xstat=$?
echo 1>&2 "$out"
echo "$out" | jq -S -r .NumChanged
exit $xstat