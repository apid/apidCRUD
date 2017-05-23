#! /bin/bash
# read from the table named file, with column line.
# Limitations:
# assumes file is not more than 1000 (maxRecs) lines.

# ----- start of mainline code
PROGDIR=$(cd "$(dirname "$0")" && /bin/pwd)
. "$PROGDIR/tester-env.sh" || exit 1
. "$PROGDIR/test-common.sh" || exit 1

TABNAME=file
KEY=line

apicurl GET "db/_table/$TABNAME?fields=$KEY" \
| jq -S -r '.records[].values[]'
