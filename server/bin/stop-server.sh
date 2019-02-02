#!/bin/bash

DIR=`readlink -f "$0"` && DIR=`dirname "$DIR"` && cd "$DIR" || exit 1

BIN="wstail-${1:-test}-server"

PID_FILE="${BIN}.pid"
LOG_FILE="${BIN}.log"

PID=''
if [ -f "$PID_FILE" ]; then
	PID=`cat "$PID_FILE"`
	EXE="`pwd`/$BIN"

	echo 'stop server' >> "$LOG_FILE" 2>&1 &

	./safe-kill.sh "$PID" "$EXE"
	echo
	while [ -e "/proc/$PID/exe" ];
	do
		sleep 1;
	done;
fi
