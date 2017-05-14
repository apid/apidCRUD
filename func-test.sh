#! /bin/bash
#	func-test.sh
# runs the server, runs tester.sh, then kills the server.
# the exit status is the exit status from tester.sh

# ----- start of mainline code
TESTS_DIR=functests
. "$TESTS_DIR/tester-env.sh" || exit 1
. "$TESTS_DIR/test-common.sh" || exit 1

LOG_DIR=${LOG_DIR:-logs}
LOGFILE=$LOG_DIR/func-test.out
NSLEEP=2

vrun ./runner.sh
echo ""

vrun ./logrun.sh "$LOGFILE" ./tester.sh
xstat=$?
sleep "$NSLEEP"
echo ""

vrun pkill -f "$DAEMON_NAME"
exit $xstat
