#! /bin/bash

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
