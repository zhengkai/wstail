package main

import (
	"fmt"
	"pb"
	"regexp"
)

var (
	comingChan    = make(chan *playerConn, 1000)
	cmdDisconnect = `DISCONNECT`

	loginPattern = regexp.MustCompile(`^[0-9a-zA-Z]{1,20}$`)
)

func manager() {

	for {
		conn := <-comingChan

		fmt.Println(`coming`, conn)

		r := getRoom(123)

		go r.playerAdd(conn)
	}
}

func loginAuth(login *pb.Login) (auth bool) {

	if !loginPattern.MatchString(login.Name) {
		return
	}
	return true
}
