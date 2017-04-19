#! /bin/bash

. ./env.sh

DAEMON=apidCRUD
MYTMP=./logs
LOGFILE=$MYTMP/$DAEMON.log
PIDFILE=$MYTMP/$DAEMON.pid
EXE=$GOPATH/bin/$DAEMON

vrun()
{
	echo 1>&2 "+ $*"
	"$@"
}

dorun()
{
	if [[ -f "$PIDFILE" ]]; then
		# vrun kill -9 "$(cat "$PIDFILE")"
		vrun pkill -f "$DAEMON"
		vrun /bin/rm -f "$PIDFILE"
	fi
	vrun "$EXE" "$@" 2>&1 &

	# this doesn't yet work correctly
	local pid=$!
	echo "$pid" > "$PIDFILE"
	wait
}

# ----- start of mainline
dorun "$@" &
sleep 2
pgrep -l $DAEMON
