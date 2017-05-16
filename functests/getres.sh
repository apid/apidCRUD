#! /bin/bash
# the API is GET on /db aka getDbResources .

# ----- start of mainline code
PROGDIR=$(cd "$(dirname "$0")" && /bin/pwd)
. "$PROGDIR/tester-env.sh" || exit 1
. "$PROGDIR/test-common.sh" || exit 1

apicurl GET "db"
