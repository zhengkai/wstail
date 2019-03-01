package main

import (
	"flag"
	"time"
)

var (
	wsWriteWait = 27 * time.Second

	gameFPS = 60

	readBuffSize = 65536

	dirBase = `/mnt/logjson/log`

	// listenHost = `127.0.0.1`
	listenHost = `0.0.0.0`
	listenPort = flag.Int("port", 21002, "http port")
)
