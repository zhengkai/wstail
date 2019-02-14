package main

import (
	"fmt"
	"log"
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
type ifConfirmChan chan *room

type room struct {
	id        int
	recv      ifRecvChan
	player    map[string]*playerConn
	isStop    bool
	idleCount int
	tickIdle  *time.Ticker
}

type roomCmd struct {
	t    roomCmdType
	conn *playerConn
	msg  *pb.Cmd
}

func getRoom(id int) (r *room) {

	load, ok := roomMap.Load(id)
	if ok {
		return load.(*room)
	}

	r = &room{
		id: id,
	}
	r.recv = make(ifRecvChan, 1000)
	roomMap.Store(id, r)
	go func() {
		r.start()
		roomMap.Delete(id)
		r.stop()
	}()

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

	r.tickIdle = time.NewTicker(5 * time.Second)

	go r.tick()

	for {
		ok := r.loopServe()
		if !ok {
			break
		}
	}
}

func (r *room) loopServe() (ok bool) {

	var recv *roomCmd

	select {
	case <-r.tickIdle.C:
		ok = r.checkIdle()
		return
	case recv, ok = <-r.recv:
	}

	if !ok {
		return
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

	return
}

func (r *room) tick() {

	t := time.NewTicker(1 * time.Second)

	var i int32

	for {

		now := <-t.C
		i++

		if r.isStop {
			break
		}

		r.recv <- &roomCmd{
			t: roomCmdBroadcast,
			msg: &pb.Cmd{
				Name: now.Format(`2006-01-02 15:04:05`),
				Test: i,
			},
		}
	}

	t.Stop()
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

func (r *room) checkIdle() (ok bool) {

	num := len(r.player)

	log.Println(`check idle`, num)

	if num == 0 {
		r.idleCount++
		if r.idleCount > 3 {
			return false
		}
	} else {
		r.idleCount = 0
	}
	return true
}

func (r *room) cmdBroadcast(c *roomCmd) {

	i := 0
	for _, v := range r.player {
		i++
		v.send <- []byte(`abc`)
	}
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
	r.isStop = true
	r.tickIdle.Stop()
	for _, v := range r.player {
		v.ws.Close()
		close(v.send)
	}
	r.player = nil
}
