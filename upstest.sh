#! /bin/bash
#	upstest.sh IDS
# IDS is a comma-separated list of record ids.
# update the given records.
# the API is PATCH on /db/_table/{table_name} aka updateDbRecords .

. tester-env.sh || exit 1
IDS=${1:-52}
RESOURCES='[{"keys":["name", "uri"], "values":["name9", "host2:z"]}]'
BODY="{\"records\":$RESOURCES}"
# echo 1>&2 "# BODY=$BODY"

out=$(./appcurl.sh PATCH "db/_table/$TABLE_NAME?ids=$IDS" -v -d "$BODY")
xstat=$?
echo 1>&2 "$out"
echo "$out" | jq -S -r .NumChanged
exit $xstat
