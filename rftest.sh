#! /bin/bash
# read from the table named file, with column line.
# Limitations:
# assumes file is not more than 1000 (maxRecs) lines.

# ----- start of mainline code
. tester-env.sh || exit 1
TABNAME=file
KEY=line

./appcurl.sh GET "db/_table/$TABNAME?fields=$KEY" \
| jq -S -r .Records[].Values[]
