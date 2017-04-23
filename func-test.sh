#! /bin/bash

vrun()
{
	echo 1>&2 "+ $*"
	"$@"
}

# ----- start of mainline code
vrun ./runner.sh
sleep 2
echo ""

vrun ./tester.sh
xstat=$?
sleep 2
echo ""

vrun pkill -f apidCRUD
exit $xstat
