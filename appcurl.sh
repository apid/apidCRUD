#! /bin/bash
#	apicurl VERB API_PATH [OTHER_CURL_ARGS]
#

vrun()
{
	echo 1>&2 "+ $*"
	"$@"
}

# strip EOL comments from the given input stream.
stripcom()
{
	sed -e 's/[ 	]*#.*//' "$@"
}

# print from the given yaml file, the value of the named top-level variable.
get_var()
{
	local yaml=$1 name=$2
	stripcom < "$yaml" \
	| sed -n -e 's/^'"$name"':[ 	]*\(.*\)/\1/p'
}

# ----- start of mainline code
CFG=./apid_config.yaml
APP=apid

API_LISTEN=$(get_var "$CFG" api_listen)

URL_BASE=http://$API_LISTEN/$APP

VERB=$1
shift
API_PATH=$1
shift
WFMT=":code:%{http_code}"

out=$(vrun curl -s -S \
	-X $VERB \
	-H "Content-type: application/json" \
	-w ":code:%{http_code}" \
	"$URL_BASE/$API_PATH" \
	"$@")
xstat=$?

# delete everything but the trailing http_code
code=${out##*:code:}

# delete the trailing marker and code
out=${out%:code:*}

echo "$out"
echo ""

if ! ((200 <= code && code < 300)); then
	xstat=1
fi
exit $xstat
