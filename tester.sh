#! /bin/bash
#	tester.sh
# try out a variety of APIs, and do some cursory tests.
# this assumes the server is already running.

get_rec_ids()
{
	./recstest.sh '*' 2>/dev/null \
	| jq -S '.Records[].Values[0]'
}

get_rec_uri()
{
	ID=$1
	TABLE=bundles
	FIELDS=uri
	API_PATH=db/_table
	VERBOSE=
	./appcurl.sh GET "$API_PATH/$TABLE/$ID?fields=$FIELDS" $VERBOSE 2>/dev/null \
	| jq -r -S '.Records[].Values[0]'
}

list_tables()
{
	echo ".tables" \
	| sqlite3 "$DBFILE" 2>/dev/null \
	| tr ' ' '\n' | grep -v '^$'
}

TestHeader()
{
	echo -n "# $* - "
}

AssertOK()
{
	if [ $? -ne 0 ]; then
		echo "FAIL - $*"
		exit 1
	else
		echo OK
	fi
}

# ----- start of mainline
. tester-env.sh || exit 1

# start clean
TestHeader creating empty database
./mkdb.sh
AssertOK "database initialization"

TestHeader checking tables "(tabtest.sh)"
out=$(./tabtest.sh 2>/dev/null | sort | tr '\n' ' ')
tabs=( $out )
exp=( bundles file nothing users )
[[ "${tabs[*]}" == "${exp[*]}" ]]
AssertOK "tabtest.sh expected [${exp[*]}], got [${tabs[*]}]"

TestHeader "adding a few records (crtest.sh)"
nrecs=7
out=$(./crtest.sh "$nrecs" 2>/dev/null | jq -S '.Ids[]')
nc=$(echo "$out" | grep -c "")
[[ "$nc" == "$nrecs" ]]
AssertOK "crtest.sh expected $nrecs, got $nc"

TestHeader "read one record (rectest.sh)"
out=$(./rectest.sh 7 2>/dev/null)
[[ "$out" == 7 ]]
AssertOK "rectest.sh expected 7, got $out"

TestHeader "reading the records (recstest.sh)"
total=$(get_rec_ids | grep -c "")
[[ "$total" == "$nc" ]]
AssertOK "recstest.sh expected $total, got $nc"

TestHeader "deleting a record (deltest.sh)"
nc=$(./deltest.sh 7 2>/dev/null)
[[ "$nc" == 1 ]]
AssertOK "deltest.sh expected 1, got $nc"

TestHeader "checking total number of records (recstest.sh)"
total=$(get_rec_ids | grep -c "")
((xtotal=nrecs-1))
[[ "$total" == "$xtotal" ]]
AssertOK "deltest.sh expected $xtotal, got $total"

TestHeader "deleting more records (delstest.sh)"
nc=$(./delstest.sh 2,3,4 2>/dev/null)
[[ "$nc" == 3 ]]
AssertOK "delstest.sh expected 3, got $nc"

TestHeader "updating a record (uptest.sh)"
nc=$(./uptest.sh 5 2>/dev/null)
[[ "$nc" == 1 ]]
AssertOK "uptest.sh expected 1, got $nc"

TestHeader "check rec 6 uri before update (get_rec_uri)"
uri1=$(get_rec_uri 6)
# echo "uri1=$uri1"
[[ $uri1 != "" ]]
AssertOK "uri1 empty"

TestHeader "updating 2 records (upstest.sh)"
nc=$(./upstest.sh 1,6 2>/dev/null)
[[ "$nc" == 2 ]]
AssertOK "upstest.sh expected 2, got $nc"

TestHeader "checking the update (get_rec_uri)"
uri2=$(get_rec_uri 6)
# echo "uri2=$uri2"
[[ "$uri1" != "$uri2" ]]
AssertOK "update did not change uri = $uri1"

TestHeader "try writing a small file and reading it back (rwftest.sh)"
./rwftest.sh cmd/apidCRUD/main.go > /dev/null 2>&1
AssertOK file comparison

TestHeader "trying table creation (crtabtest.sh)"
out=$(crtabtest.sh ABC 2>/dev/null)
out=$(list_tables | grep '^ABC$')
AssertOK "table creation"

TestHeader "trying table deletion (deltabtest.sh)"
out=$(deltabtest.sh ABC 2>/dev/null)
out=$(list_tables | grep '^$ABC$')
[[ $? != 0 ]]  # the grep should have failed
AssertOK "table deletion"

TestHeader "trying tables creation (crtabstest.sh)"
out=$(crtabstest.sh X Y Z 2>/dev/null)
out=$(list_tables | grep -c '^[XYZ]$')
[[ "$out" == 3 ]]
AssertOK "tables creation"

echo "# all passed"
exit 0
