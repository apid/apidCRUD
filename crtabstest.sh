#! /bin/bash
#	crtabstest.sh [TABNAMES]
# create multiple tables in one call.
# the API is POST /db/_schema aka createDbTables

. tester-env.sh || exit 1

tabdata()
{
	local TABNAME=$1
	local FIELD_ID='{"name":"id","properties":["primary","int32"]}'
	local FIELD_URI='{"name":"uri","properties":[]}'
	local FIELD_NAME='{"name":"name","properties":[]}'
	local FIELDS='['"$FIELD_ID,$FIELD_URI,$FIELD_NAME"']'
	local BTABLE='{"name":"'"$TABNAME"'","fields":'"$FIELDS"'}'
	echo "$BTABLE"
}

dotabs()
{
	local tsep=''
	for tab in "$@"; do
		echo "$tsep"
		tabdata "$tab"
		tsep=","
	done
}

# ----- start of mainline
BODY='{"resource":['"$(dotabs "$@")"']}'

# echo "$BODY" | jq .

if [[ $# -eq 0 ]]; then
	echo 1>&2 "error: TABLE_NAMES must be specified on cmd line"
	exit 1
fi

out=$(./appcurl.sh POST "db/_schema" -v -d "$BODY")
echo "$out"

echo "" 1>&2
echo ".tables" | sqlite3 "$DBFILE" 1>&2
