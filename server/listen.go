package main

import (
	"fmt"
	"net/http"
	"time"

	"pb"

	"github.com/gogo/protobuf/proto"
	"github.com/gorilla/websocket"
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

	closeCode := 0
	closeText := ``

	defer func() {
		if closeCode > 0 {
			ws.SetWriteDeadline(time.Now().Add(wsWriteWait))
			ab := websocket.FormatCloseMessage(closeCode, closeText)
			ws.WriteMessage(websocket.CloseMessage, ab)
		}
		ws.Close()
	}()

	ws.SetReadDeadline(time.Now().Add(30 * time.Second))

	var ab []byte
	_, ab, err = ws.ReadMessage()
	if err != nil {
		return
	}

	login := &pb.Login{}
	err = proto.Unmarshal(ab, login)

	if err != nil || login.Name == `` {
		closeCode = 1
		closeText = `login parse error`
		return
	}

	fmt.Println(`login`, login.Name)

	ticker := time.NewTicker(500 * time.Millisecond)

	i := 0

	for t := range ticker.C {

		i++
		if i > 10 {
			break
		}

		fmt.Println(t)

		if err != nil {
			break
		}
	}

	fmt.Println(`listen end`)
}
