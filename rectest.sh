#! /bin/bash
TABLE=bundles
FIELDS=${FIELDS:-id,name,uri}
API_PATH=db/_table
VERBOSE=

IDS=${1:-1,2,3}

bad=0
for i in ${IDS//,/ }; do
	out=$(./appcurl.sh GET "$API_PATH/$TABLE/$i?fields=$FIELDS&a=b&c=d" \
		$VERBOSE \
		-d "body=bodystuff")
	xstat=$?
	if [[ $xstat -ne 0 ]]; then
		bad=1
	fi
	echo 1>&2 "$out"
	echo "$out" | jq -r -S .Record[].id
done

exit $bad