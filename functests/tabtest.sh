#! /bin/bash
#	tabtest.sh
# print the names of the available tables in the database.
# the API is GET on /db/_table aka getDbTables .

# ----- start of mainline code
PROGDIR=$(cd "$(dirname "$0")" && /bin/pwd)
. "$PROGDIR/tester-env.sh" || exit 1
. "$PROGDIR/test-common.sh" || exit 1

((i=${1:-7}))
out=$(apicurl GET "db/_table")
xstat=$?
echo 1>&2 "$out"
echo "$out" | jq -S -r .Names[]
exit $xstat
