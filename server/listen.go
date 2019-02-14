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

type listen struct {
	closeCode int
	closeText string
	ws        *websocket.Conn
	login     *pb.Login
	conn      *playerConn
	room      *room
}

type playerConn struct {
	id      uint64
	login   *pb.Login
	ws      *websocket.Conn
	confirm ifConfirmChan
	send    chan []byte
}

func doListen(w http.ResponseWriter, r *http.Request) {

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	l := &listen{
		ws: ws,
	}
	l.run()
}

func (l *listen) waitLogin() (ok bool) {

	var err error
	var ab []byte

	l.ws.SetReadDeadline(time.Now().Add(21 * time.Second))
	_, ab, err = l.ws.ReadMessage()
	if err != nil {
		return
	}

	l.login = &pb.Login{}
	err = proto.Unmarshal(ab, l.login)

	if err != nil {
		l.closeCode = 1
		l.closeText = `login parse error`
		return
	}

	if !loginAuth(l.login) {
		l.closeCode = 2
		l.closeText = `login fail`
		return
	}

	return true
}

func (l *listen) waitManager() (ok bool) {

	confirmChan := make(ifConfirmChan, 1)
	l.conn = &playerConn{
		id:      atomic.AddUint64(&playerConnID, 1),
		login:   l.login,
		ws:      l.ws,
		confirm: confirmChan,
	}

	comingChan <- l.conn

	l.room = <-confirmChan

	close(confirmChan)
	l.conn.confirm = nil

	if l.room == nil {
		l.closeCode = 3
		l.closeText = `not confirmed`
		return
	}

	return true
}

func (l *listen) run() {
	for {
		if !l.waitLogin() {
			break
		}
		if !l.waitManager() {
			break
		}
		l.loopRead()
		break
	}
	l.close()
}

func (l *listen) loopRead() {

	var err error
	var ab []byte

	l.ws.SetReadDeadline(time.Time{})
	for {

		_, ab, err = l.ws.ReadMessage()
		if err != nil {
			break
		}

		fmt.Println(`client read`, string(ab))
	}

	l.room.playerExit(l.conn)
}

func (l *listen) close() {
	if l.closeCode > 0 {
		l.ws.SetWriteDeadline(time.Now().Add(wsWriteWait))
		ab := websocket.FormatCloseMessage(l.closeCode+4000, l.closeText)
		l.ws.WriteMessage(websocket.CloseMessage, ab)
	}
	l.ws.Close()
}
