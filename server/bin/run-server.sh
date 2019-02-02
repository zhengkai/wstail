#!/bin/bash

DIR=`readlink -f "$0"` && DIR=`dirname "$DIR"` && cd "$DIR" || exit 1

ulimit -n 100000
ulimit -n

BIN="wstail-${1:-test}-server"
# PORT="${2:-21009}"

./stop-server.sh "${1:-test}"

echo start build '"'$BIN'"'
TIME="time: %E" time go build -o "$BIN" ../*.go
echo

chmod +x "$BIN"

PID_FILE="${BIN}.pid"
LOG_FILE="${BIN}.log"

# nohup "./$BIN" -port "$PORT" > "$LOG_FILE" 2>&1 &
nohup "./$BIN" > "$LOG_FILE" 2>&1 &
PID="$!"
echo "$PID" > "$PID_FILE"
echo "new server started, pid = $PID"

# tail -F "$LOG_FILE"
