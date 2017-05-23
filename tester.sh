#! /bin/bash
#	tester.sh
# try out a variety of APIs, and do some cursory tests.
# this assumes the server is already running.

get_nrecs()
{
	Logrun "$TESTS_DIR/recstest.sh" '*' \
	| jq -S '.records[].values[0]' | grep -c ""
}

get_rec_uri()
{
	local ID=$1
	local TABLE=bundles
	local FIELDS=id,uri
	local API_PATH=db/_table
	Logrun apicurl GET "$API_PATH/$TABLE/$ID?fields=$FIELDS" \
	| jq -r -S '.records[].values[1]'
}

Logrun()
{
	echo "" >> "$LOG_FILE"
	echo "+ $*" >> "$LOG_FILE"
	"$@" 2>> "$LOG_FILE"
}

list_tables()
{
	echo ".tables" \
	| sqlite3 "$DBFILE" 2>/dev/null \
	| tr ' ' '\n' | grep -v '^$'
}

Logecho()
{
	echo "$@"
	echo "$@" >> "$LOG_FILE"
	if [[ "$1" == "-n" ]]; then
		echo "" >> "$LOG_FILE"
	fi
}

TestHeader()
{
	Logecho -n "# $* -"
}

AssertOK()
{
	if [ $? -ne 0 ]; then
		Logecho "FAIL - $*"
		exit 1
	else
		Logecho OK
	fi
}

# ----- start of mainline
TESTS_DIR=functests
LOG_DIR=logs
LOG_FILE=$LOG_DIR/tester.out

mkdir -p "$LOG_DIR"
/bin/rm -f "$LOG_FILE"

. "$TESTS_DIR/tester-env.sh" || exit 1
. "$TESTS_DIR/test-common.sh" || exit 1

# start clean
TestHeader creating empty database
"$TESTS_DIR/mkdb.sh"
AssertOK "database initialization"

TestHeader checking _tables_ "(tabtest.sh)"
out=$(Logrun "$TESTS_DIR/tabtest.sh" | sort | tr '\n' ' ')
tabs=( $out )
exp=( bundles file nothing users )
[[ "${tabs[*]}" == "${exp[*]}" ]]
AssertOK "tabtest.sh expected [${exp[*]}], got [${tabs[*]}]"

TestHeader "adding a few records (crtest.sh)"
nrecs=7
out=$(Logrun "$TESTS_DIR/crtest.sh" "$nrecs" | jq -S '.Ids[]')
nc=$(echo "$out" | grep -c "")
[[ "$nc" == "$nrecs" ]]
AssertOK "crtest.sh expected $nrecs, got $nc"

TestHeader "read one record (rectest.sh)"
out=$(Logrun "$TESTS_DIR/rectest.sh" 7)
[[ "$out" == 7 ]]
AssertOK "rectest.sh expected 7, got $out"

TestHeader "reading the records (recstest.sh)"
total=$(get_nrecs)
[[ "$total" == "$nc" ]]
AssertOK "recstest.sh expected $total, got $nc"

TestHeader "deleting a record (deltest.sh)"
nc=$(Logrun "$TESTS_DIR/deltest.sh" 7)
[[ "$nc" == 1 ]]
AssertOK "deltest.sh expected 1, got $nc"

TestHeader "checking total number of records (recstest.sh)"
total=$(get_nrecs)
((xtotal=nrecs-1))
[[ "$total" == "$xtotal" ]]
AssertOK "deltest.sh expected $xtotal, got $total"

TestHeader "deleting more records (delstest.sh)"
nc=$(Logrun "$TESTS_DIR/delstest.sh" 2,3,4)
[[ "$nc" == 3 ]]
AssertOK "delstest.sh expected 3, got $nc"

TestHeader "updating a record (uptest.sh)"
nc=$(Logrun "$TESTS_DIR/uptest.sh" 5)
[[ "$nc" == 1 ]]
AssertOK "uptest.sh expected 1, got $nc"

TestHeader "check rec 6 uri before update (get_rec_uri)"
uri1=$(get_rec_uri 6)
[[ $uri1 != "" ]]
AssertOK "uri1 empty"

TestHeader "updating 2 records (upstest.sh)"
nc=$(Logrun "$TESTS_DIR/upstest.sh" 1,6)
[[ "$nc" == 2 ]]
AssertOK "upstest.sh expected 2, got $nc"

TestHeader "checking the update (get_rec_uri)"
uri2=$(get_rec_uri 6)
[[ "$uri1" != "$uri2" ]]
AssertOK "update did not change uri = $uri1"

TestHeader "try writing a small file and reading it back (rwftest.sh)"
"$TESTS_DIR/rwftest.sh" cmd/apidCRUD/main.go > /dev/null 2>&1
AssertOK file comparison

TestHeader "trying tables creation (crtabtest.sh)"
out=$(Logrun "$TESTS_DIR/crtabtest.sh" X Y Z)
out=$(list_tables | grep -c '^[XYZ]$')
[[ "$out" == 3 ]]
AssertOK "tables creation"

TestHeader "trying table deletion (deltabtest.sh)"
out=$(Logrun "$TESTS_DIR/deltabtest.sh" X Y Z)
out=$(list_tables | grep '^$[XYZ]$')
[[ $? != 0 ]]  # the grep should have failed
AssertOK "table deletion"

TestHeader "trying table schema (desctabtest.sh)"
out=$(Logrun "$TESTS_DIR/desctabtest.sh" users)
[[ "$out" != "" ]]
AssertOK "table description"

TestHeader "trying db resources (getres.sh)"
out=$(Logrun "$TESTS_DIR/getres.sh" h)
xstat=$?
[[ $xstat == 0 && "$out" != "" ]]
AssertOK "db resources"

Logecho "# all passed"
exit 0
