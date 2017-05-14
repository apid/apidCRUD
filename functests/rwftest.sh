#! /bin/bash

notice()
{
	echo 1>&2 "# $*"
}

# ----- start of mainline code
PROGDIR=$(cd "$(dirname "$0")" && /bin/pwd)
. "$PROGDIR/tester-env.sh" || exit 1
. "$PROGDIR/test-common.sh" || exit 1

TESTFILE=${1:-main}
PROGNAME=${0##*/}
TMPFILE=/tmp/$PROGNAME-$$.tmp

notice 'creating empty table "file"'
echo "drop table if exists file;create table file(line text);" \
| sqlite3 "$DBFILE"

notice "copying $TESTFILE" 'to table "file"'
"$TESTS_DIR/wftest.sh" "$TESTFILE"

notice "reading back $TESTFILE" 'from table "file"'
"$TESTS_DIR/rftest.sh" > "$TMPFILE"

notice "diffing the result"
diff "$TESTFILE" "$TMPFILE"
xstat=$?
/bin/rm -f "$TMPFILE"

if [[ $xstat -ne 0 ]]; then
	notice "FAIL - result is not the same"
	exit 1
fi
echo OK
