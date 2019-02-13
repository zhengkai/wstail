package main

import (
	"time"
)

func main() {

	go manager()
	go server()
	time.Sleep(time.Hour * 999999)
}
