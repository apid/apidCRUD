#! /bin/bash
#	deltabtest.sh [TABNAME]
# delete a table.
# the API is DELETE /db/_schema aka createDbTable

. tester-env.sh || exit 1
TABNAME=${1:-anewtab}

out=$(./appcurl.sh DELETE "db/_schema/$TABNAME" -v -d "$BODY")
echo "$out"

echo 1>&2 ""
echo ".tables" | sqlite3 "$DBFILE" 1>&2
