#! /bin/bash
#	appcurl VERB API_PATH [OTHER_CURL_ARGS]
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
get_config_var()
{
	local yaml=$1 name=$2
	stripcom < "$yaml" \
	| sed -n -e 's/^'"$name"':[ 	]*\(.*\)/\1/p'
}

apicurl()
{
	local VERB=$1 API_PATH=$2
	shift 2

	local WFMT=":code:%{http_code}"
	local out
	out=$(vrun curl -s -S \
		-X "$VERB" \
		-H "Content-type: application/json" \
		-w ":code:%{http_code}" \
		"$URL_BASE/$API_PATH" \
		"$@")
	local xstat=$?

	# delete everything but the trailing http_code
	local code=${out##*:code:}

	# delete the trailing marker and code
	local out=${out%:code:*}

	if [[ -n "$out" ]]; then
		echo "$out"
		echo ""
	fi

	if ! ((200 <= code && code < 300)); then
		xstat=11
	fi
	return $xstat
}

# ----- start of mainline code
CFG=./apid_config.yaml
API_PREFIX=apid

API_LISTEN=$(get_config_var "$CFG" api_listen)
URL_BASE=http://$API_LISTEN/$API_PREFIX

apicurl "$@"
