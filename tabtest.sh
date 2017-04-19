#! /bin/bash
#	tabtest.sh
# print the names of the available tables in the database.

FIELDS=id,name,uri
API_PATH=db/_table

((i=${1:-7}))
out=$(./appcurl.sh GET "$API_PATH")
xstat=$?
echo 1>&2 "$out"
echo "$out" | jq -S -r .Resource[]
exit $xstat
