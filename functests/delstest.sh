#! /bin/bash
#	delstest.sh
# functional test for the DELETE records API.
# the API is DELETE /db/_table/{table_name} aka deleteDbRecords .

# ----- start of mainline code
PROGDIR=$(cd "$(dirname "$0")" && /bin/pwd)
. "$PROGDIR/tester-env.sh" || exit 1
. "$PROGDIR/test-common.sh" || exit 1

IDS=${1:-25}
ID_FIELD=id
out=$(apicurl DELETE \
	"db/_table/$TABLE_NAME?ids=$IDS&id_field=$ID_FIELD")

echo 1>&2 "$out"
echo "$out" | jq -S -r .numChanged
