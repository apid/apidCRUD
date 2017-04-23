#! /bin/bash
#	logrun LOGFILE CMD ...
# run the given CMD with stdout/stderr going to LOGFILE.
# exit with the status of CMD.

LOGFILE=$1; shift

"$@" | tee "$LOGFILE"
exit ${PIPESTATUS[0]}
