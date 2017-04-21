#! /bin/bash
#	tester.sh
# try out a variety of APIs, and do some cursory tests.
# this assumes the server is already running.

get_rec_ids()
{
	./recstest.sh '*' 2>/dev/null \
	| jq -S '.Record' \
	| jq -S -r '.[].id'
}

fail()
{
	echo "FAIL - $*"
	exit 1
}

# ----- start of mainline
# start clean
echo "# creating empty database"
mkdb.sh || exit 1
echo OK

echo "# checking tables"
out=$(./tabtest.sh 2>/dev/null | sort | tr '\n' ' ')
tabs=( $out )
exp=( bundles nothing users )
if [[ "${tabs[*]}" != "${exp[*]}" ]]; then
	fail "tabtest.sh expected [${exp[*]}], got [${tabs[*]}]"
else
	echo OK
fi


echo "# adding a few records"
nrecs=7
out=$(./crtest.sh "$nrecs" 2>/dev/null | jq -S '.Ids[]')
nc=$(echo "$out" | wc -l)
if [[ "$nc" != "$nrecs" ]]; then
	fail "crtest.sh expected $nrecs, got $nc"
else
	echo OK
fi

echo "# read one record"
out=$(./rectest.sh 7 2>/dev/null)
if [[ "$out" != 7 ]]; then
	fail "rectest.sh expected 7, got $out"
else
	echo OK
fi

echo "# reading the records"
total=$(get_rec_ids | wc -l)
if [[ "$total" != "$nc" ]]; then
	fail "recstest.sh expected $total, got $nc"
else
	echo OK
fi

echo "# deleting a record"
nc=$(./deltest.sh 7 2>/dev/null)
if [[ "$nc" != 1 ]]; then
	fail "deltest.sh expected 1, got $nc"
else
	echo OK
fi

echo "# checking total number of records"
total=$(get_rec_ids | wc -l)
((xtotal=nrecs-1))
if [[ "$total" != "$xtotal" ]]; then
	fail "deltest.sh expected $xtotal, got $total"
else
	echo OK
fi

echo "# deleting more records"
nc=$(./delstest.sh 2,3,4 2>/dev/null)
if [[ "$nc" != 3 ]]; then
	fail "delstest.sh expected 3, got $nc"
else
	echo OK
fi

echo "# updating a record"
nc=$(./uptest.sh 5 2>/dev/null)
if [[ "$nc" != 1 ]]; then
	fail "uptest.sh expected 1, got $nc"
else
	echo OK
fi

echo "# updating 2 records"
nc=$(./upstest.sh 1,6 2>/dev/null)
if [[ "$nc" != 2 ]]; then
	fail "upstest.sh expected 2, got $nc"
else
	echo OK
fi

echo "# all good"
exit 0
