#! /bin/bash
#	func-test.sh
# runs the server, runs tester.sh, then kills the server.
# the exit status is the exit status from tester.sh

vrun()
{
	echo 1>&2 "+ $*"
	"$@"
}

# ----- start of mainline code
NSLEEP=2

vrun ./runner.sh
sleep "$NSLEEP"
echo ""

vrun ./tester.sh
xstat=$?
sleep "$NSLEEP"
echo ""

vrun pkill -f apidCRUD
exit $xstat
