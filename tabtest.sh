#! /bin/bash
#	tabtest.sh
# print the names of the available tables in the database.
# the API is GET on /db/_table aka getDbTables .

FIELDS=id,name,uri
API_PATH=db/_table

((i=${1:-7}))
out=$(./appcurl.sh GET "$API_PATH")
xstat=$?
echo 1>&2 "$out"
echo "$out" | jq -S -r .Names[]
exit $xstat
