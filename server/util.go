package main

import (
	"bufio"
	"os"
	"time"

	"pb"
)

func createReader(filename string) (r *bufio.Reader, err error) {
	f, err := os.Open(filename)
	if err != nil {
		return
	}
	r = bufio.NewReader(f)
	return
}

func makeOpBaseReturn() (r *pb.OpBaseReturn) {

	ts := time.Now().UnixNano() / 1000000

	r = &pb.OpBaseReturn{
		ServerTs: uint64(ts),
	}

	return
}
