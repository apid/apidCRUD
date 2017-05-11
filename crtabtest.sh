#! /bin/bash
#	crtabtest.sh [TABNAME]
# create a table.
# the API is POST /db/_schema aka createDbTable

DBFILE=apidCRUD.db

TABNAME=${1:-anewtab}
FIELD_ID='{"name":"id","properties":["primary","int32"]}'
FIELD_URI='{"name":"uri","properties":[]}'
FIELD_NAME='{"name":"name","properties":[]}'
FIELDS='['"$FIELD_ID,$FIELD_URI,$FIELD_NAME"']'
BTABLE='{"name":"'"$TABNAME"'","fields":'"$FIELDS"'}'
TABLES='['"$BTABLE"']'
BODY='{"resource":'"$TABLES"'}'

echo "$BODY" | jq .

out=$(./appcurl.sh POST "db/_schema/$TABNAME" -v -d "$BODY")
echo "$out"

echo "" 1>&2
echo ".tables" | sqlite3 "$DBFILE" 1>&2
