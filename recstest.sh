#! /bin/bash
#	recstest.sh IDS
# dump all records (by default)
# the API is GET on /db/_table/{table_name} aka getDbRecords

. tester-env.sh || exit 1
FIELDS=id,name
API_PATH=db/_table
IDS=${1:-\*}
if [[ "$IDS" == \* ]]; then
	IDS=
fi

out=$(./appcurl.sh GET "$API_PATH/$TABLE_NAME?ids=$IDS&fields=$FIELDS")
xstat=$?

echo "$out"
# echo "$out" | jq -r -S .Records[].id
exit $xstat
