#! /bin/bash
# write the given file to the table "file".
# Limitations:
# the file should be not more than 1000 (maxRecs) lines.
# may not contain arbitrary control chars.

#
# escape some of the common characters that JSON
# requires to be escaped in a string.
# this should probably be done in perl.
#
jescape()
{
	sed -e 's/["\\]/\\&/g' \
		-e 's/	/\\t/g' \
		-e 's//\\r/g' \
		"$@"
}

mkrecs()
{
	local line
	local sep=''
	while IFS= read -r line; do
		cat<<EOF
$sep{"keys":["line"], "values":["$line"]}
EOF
		sep=,
	done
}

# ----- start of mainline code
FILE=${1:-swagger.yaml}
TABLE=file
RESOURCES="[$(jescape "$FILE" | mkrecs)]"
BODY="{\"records\":$RESOURCES}"
# echo 1>&2 "BODY=$BODY"

./appcurl.sh POST "db/_table/$TABLE" -d "$BODY"
