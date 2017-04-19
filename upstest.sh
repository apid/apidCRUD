#! /bin/bash

IDS=${1:-52}
TABLE=bundles
RESOURCES='[{"keys":["name", "uri"], "values":["name9", "host2:z"]}]'
BODY="{\"resource\":$RESOURCES}"
# echo 1>&2 "# BODY=$BODY"

out=$(./appcurl.sh PATCH "db/_table/$TABLE?ids=$IDS" -v -d "$BODY")
xstat=$?
echo 1>&2 "$out"
echo "$out" | jq -S -r .NumChanged
exit $xstat
