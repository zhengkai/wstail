package main

import (
	"fmt"
	"pb"

	"github.com/gorilla/websocket"
)

var (
	comingChan    = make(chan *playerConn, 1000)
	cmdDisconnect = `DISCONNECT`
)

type ifConfirmChan chan *room

type playerConn struct {
	id      uint64
	login   *pb.Login
	ws      *websocket.Conn
	confirm ifConfirmChan
	send    chan []byte
}

func manager() {

	for {
		conn := <-comingChan

		fmt.Println(`coming`, conn)

		r := getRoom(123)

		go r.playerAdd(conn)
	}
}

func loginAuth(login *pb.Login) (auth bool) {
	if login.Name == `` {
		return
	}
	return true
}
