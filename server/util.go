package main

import (
	"bufio"
	"os"
)

func createReader(filename string) (r *bufio.Reader, err error) {
	f, err := os.Open(filename)
	if err != nil {
		return
	}
	r = bufio.NewReader(f)
	return
}
