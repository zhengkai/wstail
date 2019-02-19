package main

import (
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/zhengkai/rome"
)

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  4096,
		WriteBufferSize: 4096,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

func doListen(w http.ResponseWriter, r *http.Request) {

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	p := &player{}
	p.WS = ws

	rome.ParsePlayerConn(p, listenChan)
}
