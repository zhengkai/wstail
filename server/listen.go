package main

import (
	"fmt"
	"net/http"
	"sync/atomic"
	"time"

	"pb"

	"github.com/gogo/protobuf/proto"
	"github.com/gorilla/websocket"
)

var (
	playerConnID uint64

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

	closeCode := 0
	closeText := ``

	defer func() {
		if closeCode > 0 {
			ws.SetWriteDeadline(time.Now().Add(wsWriteWait))
			ab := websocket.FormatCloseMessage(closeCode+4000, closeText)
			ws.WriteMessage(websocket.CloseMessage, ab)
		}
		ws.Close()
	}()

	var ab []byte
	ws.SetReadDeadline(time.Now().Add(21 * time.Second))
	_, ab, err = ws.ReadMessage()
	if err != nil {
		return
	}

	login := &pb.Login{}
	err = proto.Unmarshal(ab, login)

	if err != nil {
		closeCode = 1
		closeText = `login parse error`
		return
	}

	if !loginAuth(login) {
		closeCode = 2
		closeText = `login fail`
		return
	}

	confirmChan := make(ifConfirmChan, 1)
	conn := &playerConn{
		id:      atomic.AddUint64(&playerConnID, 1),
		login:   login,
		ws:      ws,
		confirm: confirmChan,
	}

	comingChan <- conn

	room := <-confirmChan

	close(confirmChan)
	conn.confirm = nil

	if room == nil {
		closeCode = 3
		closeText = `not confirmed`
		return
	}

	ws.SetReadDeadline(time.Time{})
	for {

		_, ab, err = ws.ReadMessage()
		if err != nil {
			break
		}

		fmt.Println(`client read`, string(ab))

		/*
			cmd := &pb.Cmd{}
			err = proto.Unmarshal(ab, cmd)
			if err != nil {
				break
			}

			room.playerMsg(conn, cmd)
		*/
	}

	room.playerExit(conn)
}
