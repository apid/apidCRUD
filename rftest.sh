#! /bin/bash
# read from the table named file, with column line.
# Limitations:
# assumes file is not more than 1000 (maxRecs) lines.

# ----- start of mainline code
FILE=${1:-swagger.yaml}
TABLE=file
KEY=line

./appcurl.sh GET "db/_table/$TABLE?fields=$KEY" -v \
| jq -S -r .Records[].Values[]
