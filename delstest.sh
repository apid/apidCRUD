#! /bin/bash

IDS=${1:-25}
TABLE=bundles
ID_FIELD=id
out=$(./appcurl.sh DELETE \
	"db/_table/$TABLE?ids=$IDS&id_field=$ID_FIELD")

echo 1>&2 "$out"
echo "$out" | jq -S -r .NumChanged
