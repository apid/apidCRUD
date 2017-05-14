#! /bin/bash
#	deltabtest.sh TABNAMES
# delete a table.
# the API is DELETE /db/_schema/XXX aka createDbTable

# ----- start of mainline code
PROGDIR=$(cd "$(dirname "$0")" && /bin/pwd)
. "$PROGDIR/tester-env.sh" || exit 1
. "$PROGDIR/test-common.sh" || exit 1

if [[ $# -eq 0 ]]; then
	echo 1>&2 "error: TABNAME must be specified on cmd line"
	exit 1
fi

for tab in "$@"; do
out=$(apicurl DELETE "db/_schema/$tab" -v)
echo "$out"
done

echo 1>&2 ""
echo ".tables" | sqlite3 "$DBFILE" 1>&2
