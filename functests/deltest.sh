#! /bin/bash
#	deltest.sh
# functional test for the DELETE record API.
# the API is DELETE on /db/_table/{table_name}/{id} aka deleteDbRecord .

# ----- start of mainline code
PROGDIR=$(cd "$(dirname "$0")" && /bin/pwd)
. "$PROGDIR/tester-env.sh" || exit 1
. "$PROGDIR/test-common.sh" || exit 1

ID=${1:-25}
out=$(apicurl DELETE "db/_table/$TABLE_NAME/$ID")
xstat=$?
echo 1>&2 "$out"
echo "$out" | jq -S -r .numChanged
exit $xstat
