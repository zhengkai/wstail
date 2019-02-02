#!/bin/bash

DIR=`readlink -f "$0"` && DIR=`dirname "$DIR"` && cd "$DIR" || exit 1

BIN="wstail-${1:-test}-server"

PID_FILE="${BIN}.pid"

PID=''
if [ -f "$PID_FILE" ]; then
	PID=`cat "$PID_FILE"`
	EXE="`pwd`/$BIN"

	./safe-status.sh "$PID" "$EXE"
fi
