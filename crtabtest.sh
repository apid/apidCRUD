#! /bin/bash
#	crtabtest.sh
# create a table.
# the API is POST /db/_schema aka createDbTable

DBFILE=apidCRUD.db

TABLE=newtab
FIELD_ID='{"name":"id","properties":["primary","int32"]}'
FIELD_URI='{"name":"uri","properties":[]}'
FIELD_NAME='{"name":"name","properties":[]}'
FIELDS='['"$FIELD_ID,$FIELD_URI,$FIELD_NAME"']'
BTABLE='{"name":"'"$TABLE"'","fields":'"$FIELDS"'}'
TABLES='['"$BTABLE"']'
BODY='{"resource":'"$TABLES"'}'

echo "$BODY" | jq .

out=$(./appcurl.sh POST "db/_schema/$TABLE" -v -d "$BODY")
echo "$out"

echo "" 1>&2
echo ".tables" | sqlite3 "$DBFILE" 1>&2
