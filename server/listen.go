package main

import (
	"fmt"
	"net/http"
	"time"

	"pb"

	"github.com/gogo/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"
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

	ticker := time.NewTicker(500 * time.Millisecond)

	i := 0

	for t := range ticker.C {

		i++
		if i > 10 {
			break
		}

		a := &pb.Play{
			Foo: `foo`,
			Bar: uint32(i),
		}

		base := &pb.OpBaseReturn{
			ServerTs: uint64(t.Unix()),
			Error:    `没错`,
		}

		c := &pb.MsgA{
			Base: base,
			Msg:  make([]*any.Any, 0),
		}

		auth := &pb.GameAuth{
			Id:   1,
			Sign: `auth`,
		}

		tmp, _ := ptypes.MarshalAny(a)
		c.Msg = append(c.Msg, tmp)

		tmp, _ = ptypes.MarshalAny(auth)
		c.Msg = append(c.Msg, tmp)

		b, _ := proto.Marshal(c)

		ws.SetWriteDeadline(time.Now().Add(wsWriteWait))
		err = ws.WriteMessage(websocket.BinaryMessage, b)

		x := &pb.MsgA{}

		proto.Unmarshal(b, x)
		fmt.Println(x, x.Base.Error)

		if err != nil {
			break
		}
	}

	fmt.Println(`listen end`)

	ws.Close()
}
