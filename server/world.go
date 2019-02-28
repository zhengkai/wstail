package main

import (
	"bufio"
	"fmt"
	"time"

	"pb"

	"github.com/zhengkai/rome"
)

type world struct {
	rome.World
	file *file

	fileSize int64

	buf      []byte
	prevLine []byte
	reader   *bufio.Reader

	online     int
	timeEmpty  time.Time
	timeCreate time.Time
}

func (w *world) init(filename string) {

	fmt.Println(`world init`, filename)

	w.file = &file{
		filename: filename,
	}
	w.file.start()

	now := time.Now()
	w.timeEmpty = now
	w.timeCreate = now

	r := &room{}
	r.World = w

	w.Room = r

	rome.InitRoom(r)
}

func (w *world) Tick(i int) (ok bool) {

	ok = true

	// w.echoTick(i)

	if !w.isChange() {
		return
	}

	w.scan()

	return
}

func (w *world) echoTick(i int) {

	if i%gameFPS == 0 {
		fmt.Println(`world tick`, w.fileSize)
		w.Room.SendMsg([]byte(fmt.Sprintf(`server tick %d`, i)))
	}
}

func (w *world) isChange() (ok bool) {

	select {
	case <-w.file.ch:
		ok = true
	default:
	}

	return
}

func (w *world) scan() {

	b, reset, ok := w.file.scan()
	if !ok {
		return
	}

	send := encodeReturn(&pb.Update{
		Reset_: reset,
		Msg:    utf8string(string(b)),
	})

	w.Room.SendMsg(send)
}

func (w *world) Player(p rome.IPlayerConn, status bool) {

	w.playerCount(status)

	if !status {
		return
	}

	go func() {

		b, ok := w.file.prev()
		if !ok {
			return
		}

		send := encodeReturn(&pb.PrevContent{
			Msg: utf8string(string(b)),
		})

		p.Send(send)
	}()

	// player := p.(*player)
}

func (w *world) Input(p rome.IPlayerConn, msg interface{}) {

	// player := p.(*player)
}

func (w *world) stop() {

	fmt.Println(`world close`, w.file.filename, time.Now().Sub(w.timeCreate))

	w.file.stop = true

	w.Room.Stop()
}

func (w *world) playerCount(status bool) {

	add := 1
	if !status {
		add = -1
	}
	w.online += add

	if w.online < 0 {
		panic(`player count -1`)
	}

	if w.online == 0 {
		w.timeEmpty = time.Now()
	}
}
