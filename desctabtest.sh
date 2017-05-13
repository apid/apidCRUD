#! /bin/bash
#	desctabtest.sh TABNAMES
# describe a table.
# the API is GET /db/_schema/XXX aka describeDbTable

. tester-env.sh || exit 1
if [[ $# -eq 0 ]]; then
	echo 1>&2 "error: TABNAME must be specified on cmd line"
	exit 1
fi

for tab in "$@"; do
out=$(./appcurl.sh GET "db/_schema/$tab" -v)
echo "$out" | jq -r -S .Schema
done
