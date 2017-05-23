#! /bin/bash
# the API is GET on /db aka describeService .

# ----- start of mainline code
PROGDIR=$(cd "$(dirname "$0")" && /bin/pwd)
. "$PROGDIR/tester-env.sh" || exit 1
. "$PROGDIR/test-common.sh" || exit 1

out=$(apicurl GET "db")
echo 1>&2 "$out"
echo "$out" | jq -S -r '.description'
