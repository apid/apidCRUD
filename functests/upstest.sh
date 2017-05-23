#! /bin/bash
#	upstest.sh IDS
# IDS is a comma-separated list of record ids.
# update the given records.
# the API is PATCH on /db/_table/{table_name} aka updateDbRecords .

# ----- start of mainline code
PROGDIR=$(cd "$(dirname "$0")" && /bin/pwd)
. "$PROGDIR/tester-env.sh" || exit 1
. "$PROGDIR/test-common.sh" || exit 1

IDS=${1:-52}
RESOURCES='[{"keys":["name", "uri"], "values":["name9", "host2:z"]}]'
BODY="{\"records\":$RESOURCES}"
# echo 1>&2 "# BODY=$BODY"

out=$(apicurl PATCH "db/_table/$TABLE_NAME?ids=$IDS" -v -d "$BODY")
xstat=$?
echo 1>&2 "$out"
echo "$out" | jq -S -r .numChanged
exit $xstat
