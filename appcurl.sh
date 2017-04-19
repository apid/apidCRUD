#! /bin/bash
#	apicurl VERB API_PATH [OTHER_CURL_ARGS]
#

vrun()
{
	echo 1>&2 "+ $*"
	"$@"
}

# ----- start of mainline code
TESTHOST=localhost
TESTPORT=9000
APP=apid
URL_BASE=http://$TESTHOST:$TESTPORT/$APP

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
