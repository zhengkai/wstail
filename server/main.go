package main

import "time"

func main() {
	go server()
	time.Sleep(time.Hour * 999999)
}
