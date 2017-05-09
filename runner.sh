#! /bin/bash

. ./env.sh

DAEMON=apidCRUD
LOG_DIR=${LOG_DIR:-./logs}
LOGFILE=$LOG_DIR/$DAEMON.log
EXE=$GOPATH/bin/$DAEMON

vrun()
{
	echo 1>&2 "+ $*"
	"$@"
}

dorun()
{
	vrun pkill -f "$DAEMON"
	# vrun "$EXE" "$@" > "$LOGFILE" 2>&1 &
	vrun "$EXE" "$@" 2>&1 &
}

# ----- start of mainline
mkdir -p "$(dirname "$LOGFILE")"
NSLEEP=2
dorun "$@" &

sleep "$NSLEEP"

# cause a failure exit if the daemon didn't start
pgrep -l $DAEMON
