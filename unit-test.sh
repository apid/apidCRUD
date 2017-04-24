#! /bin/bash
#	unit-test.sh
# run all unit tests, then generate coverage reports.

vrun()
{
	echo 1>&2 "+ $*"
	"$@"
}

# ----- start of mainline code
LOG_DIR=${LOG_DIR:-logs}

COV_DIR=${COV_DIR:-cov}
COV_FILE=${COV_FILE:-$COV_DIR/covdata.out}
COV_HTML=${COV_HTML:-$COV_DIR/apidCRUD-coverage.html}

./logrun.sh "$LOG_DIR/unit-test.out" \
go test -coverprofile="$COV_FILE"

go tool cover -func="$COV_FILE" > "$LOG_DIR/cover-func.out"

#
go tool cover -html="$COV_FILE" -o "$COV_HTML"

./tested_funcs.sh | sort > "$LOG_DIR/covered.out"
./uncovered.sh > "$LOG_DIR/uncovered.out"

