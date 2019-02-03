package main

import (
	"fmt"
	"net/http"
	"time"

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

	fmt.Println(`listen`)

	t := time.NewTicker(500 * time.Millisecond)

	i := 0

	for range t.C {

		i++
		if i > 10 {
			break
		}

		ws.SetWriteDeadline(time.Now().Add(wsWriteWait))
		err = ws.WriteMessage(websocket.TextMessage, []byte(`tick`))

		if err != nil {
			break
		}
	}

	fmt.Println(`listen end`)

	ws.Close()
}
