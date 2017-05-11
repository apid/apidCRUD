#! /bin/bash
#	deltest.sh
# functional test for the DELETE record API.
# the API is DELETE on /db/_table/{table_name}/{id} aka deleteDbRecord .

. tester-env.sh || exit 1
ID=${1:-25}
out=$(./appcurl.sh DELETE "db/_table/$TABLE_NAME/$ID")
xstat=$?
echo 1>&2 "$out"
echo "$out" | jq -S -r .NumChanged
exit $xstat
