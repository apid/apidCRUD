#! /bin/bash
#	uptest.sh ID
# update the record of the given ID.
# the API is PATCH on /db/_table/{table_name}/{id} aka updateDbRecord .

. tester-env.sh || exit 1
ID=${1:-2}
RESOURCES='[{"keys":["name", "uri"], "values":["name9", "host2:z"]}]'
BODY="{\"records\":$RESOURCES}"
# echo 1>&2 "# BODY=$BODY"

out=$(./appcurl.sh PATCH "db/_table/$TABLE_NAME/$ID" -v -d "$BODY")
xstat=$?
echo 1>&2 "$out"
echo "$out" | jq -S -r .NumChanged
exit $xstat
