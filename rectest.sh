#! /bin/bash
# the API is GET on /db/_table/{table_name}/{id} aka getDbRecord .

. tester-env.sh || exit 1
FIELDS=${FIELDS:-id,name,uri}
API_PATH=db/_table
VERBOSE=

IDS=${1:-1,2,3}

bad=0
for i in ${IDS//,/ }; do
	out=$(./appcurl.sh GET "$API_PATH/$TABLE_NAME/$i?fields=$FIELDS&a=b&c=d" \
		$VERBOSE)
	xstat=$?
	if [[ $xstat -ne 0 ]]; then
		bad=1
	fi
	echo 1>&2 "$out"
	echo "$out" | jq -r -S .Records[].Values[0]
done

exit $bad
