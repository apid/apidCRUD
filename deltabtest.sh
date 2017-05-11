#! /bin/bash
#	deltabtest.sh [TABNAME]
# delete a table.
# the API is DELETE /db/_schema aka createDbTable

DBFILE=apidCRUD.db
TABNAME=${1:-anewtab}

out=$(./appcurl.sh DELETE "db/_schema/$TABNAME" -v -d "$BODY")
echo "$out"

echo 1>&2 ""
echo ".tables" | sqlite3 "$DBFILE" 1>&2
