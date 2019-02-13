package main

import (
	"fmt"
	"pb"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type roomCmdType uint8

var (
	roomMap = sync.Map{}
)

const (
	_ roomCmdType = iota
	roomCmdMsg
	roomCmdNew
	roomCmdExit
	roomCmdBroadcast
)

type ifRecvChan chan *roomCmd

type room struct {
	id     int
	recv   ifRecvChan
	player map[string]*playerConn
	isStop bool
}

type roomCmd struct {
	t    roomCmdType
	conn *playerConn
	msg  *pb.Cmd
}

func getRoom(rID int) (r *room) {

	load, ok := roomMap.Load(rID)
	if ok {
		return load.(*room)
	}

	r = &room{
		id: rID,
	}
	r.recv = make(ifRecvChan, 1000)
	roomMap.Store(rID, r)
	go r.start()

	return
}

func (r *room) playerAdd(p *playerConn) {

	r.recv <- &roomCmd{
		t:    roomCmdNew,
		conn: p,
	}

	p.confirm <- r
}

func (r *room) playerMsg(p *playerConn, msg *pb.Cmd) {
	r.recv <- &roomCmd{
		t:    roomCmdMsg,
		conn: p,
	}
}

func (r *room) playerExit(p *playerConn) {
	r.recv <- &roomCmd{
		t:    roomCmdExit,
		conn: p,
	}
}

func (r *room) start() {

	fmt.Println(`room start`)

	r.player = make(map[string]*playerConn)

	go r.tick()

	for {

		recv, ok := <-r.recv
		if !ok {
			break
		}

		switch recv.t {

		case roomCmdNew:
			r.cmdNew(recv)

		case roomCmdExit:
			r.cmdExit(recv)

		case roomCmdMsg:
			r.cmdMsg(recv)

		case roomCmdBroadcast:
			r.cmdBroadcast(recv)
		}
	}
}

func (r *room) tick() {

	t := time.Tick(time.Second)

	var i int32

	for {

		now := <-t
		i++

		if r.isStop {
			break
		}

		fmt.Println(`tick`, i)

		r.recv <- &roomCmd{
			t: roomCmdBroadcast,
			msg: &pb.Cmd{
				Name: now.Format(`2006-01-02 15:04:05`),
				Test: i,
			},
		}
	}
}

func (r *room) cmdNew(c *roomCmd) {

	name := c.conn.login.Name

	old, ok := r.player[name]
	if ok {
		close(old.send)
		old.ws.Close()
	}

	if r.isStop {
		delete(r.player, name)
		return
	}

	c.conn.send = make(chan []byte, 1000)
	r.player[name] = c.conn

	go sendPlayer(c.conn)
}

func sendPlayer(p *playerConn) {

	for {
		ab, ok := <-p.send
		if !ok {
			return
		}

		p.ws.SetWriteDeadline(time.Now().Add(wsWriteWait))
		err := p.ws.WriteMessage(websocket.BinaryMessage, ab)
		if err != nil {
			break
		}
	}

	p.ws.Close()
}

func (r *room) cmdBroadcast(c *roomCmd) {

	i := 0
	for _, v := range r.player {
		i++
		v.send <- []byte(`abc`)
	}
	fmt.Println(`cast`, i)
}

func (r *room) cmdExit(c *roomCmd) {

	name := c.conn.login.Name

	p, ok := r.player[name]
	if !ok {
		return
	}

	if p.id != c.conn.id {
		return
	}

	delete(r.player, name)

	close(p.send)
	p.ws.Close()
}

func (r *room) cmdMsg(c *roomCmd) {
	if r.isStop {
		return
	}
}

func (r *room) stop() {
	roomMap.Delete(r.id)
	r.isStop = true
}
