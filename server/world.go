package main

import (
	"bufio"
	"fmt"

	"pb"

	"github.com/zhengkai/rome"
)

type world struct {
	rome.World
	id   int
	file *file

	fileSize int64

	buf      []byte
	prevLine []byte
	reader   *bufio.Reader
}

func (w *world) init(id int) {

	w.id = id
	w.file = filePool.get(`/tmp/fortune.txt`)
	fmt.Println(`file ok`)

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

	fmt.Printf("player type: %T\n", p)

	player := p.(*player)

	fmt.Println(`world Player`, player.ID, status)
}

func (w *world) Input(p rome.IPlayerConn, msg interface{}) {

	player := p.(*player)

	// fmt.Printf("player type: %T\n", p)
	fmt.Println(`world Input`, player.ID, msg)

	// p.SendStop([]byte(`close`))
}
