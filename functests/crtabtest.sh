#! /bin/bash
#	crtabtest.sh TABNAMES
# create multiple tables in one call.
# the API is POST /db/_schema/XXX aka createDbTables

dotab()
{
	local TABNAME=$1
	local FIELD_ID='{"name":"id","properties":["is_primary_key","int32"]}'
	local FIELD_URI='{"name":"uri","properties":[]}'
	local FIELD_NAME='{"name":"name","properties":[]}'
	local FIELDS='['"$FIELD_ID,$FIELD_URI,$FIELD_NAME"']'
	local BTABLE='{"fields":'"$FIELDS"'}'

	# echo "$BTABLE" | jq .

	apicurl POST "db/_schema/$TABNAME" -v -d "$BTABLE"
}

# ----- start of mainline code
PROGDIR=$(cd "$(dirname "$0")" && /bin/pwd)
. "$PROGDIR/tester-env.sh" || exit 1
. "$PROGDIR/test-common.sh" || exit 1

if [[ $# -eq 0 ]]; then
	echo 1>&2 "error: TABNAMES must be specified on cmd line"
	exit 1
fi

for tab in "$@"; do
	dotab "$tab"
done
echo "" 1>&2
echo ".tables" | sqlite3 "$DBFILE" 1>&2
