#! /bin/bash
#	deltabtest.sh
# delete a table.
# the API is DELETE /db/_schema aka createDbTable

DBFILE=apidCRUD.db
TABLE=newtab

out=$(./appcurl.sh DELETE "db/_schema/$TABLE" -v -d "$BODY")
echo "$out"

echo 1>&2 ""
echo ".tables" | sqlite3 "$DBFILE" 1>&2
