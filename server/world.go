package main

import (
	"bufio"
	"fmt"
	"io"
	"os"

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

	finfo, err := os.Stat(w.file.filename)
	if err != nil {
		w.reader = nil
		return
	}

	size := finfo.Size()
	fmt.Println(`size`, size, w.fileSize)

	if w.fileSize > size {
		w.reader = nil
	}
	if w.reader == nil {
		f, err := os.Open(w.file.filename)
		if err != nil {
			return
		}
		w.reader = bufio.NewReader(f)
		w.fileSize = 0
	}

	n, err := w.reader.Read(w.buf)
	if err != nil && err != io.EOF {
		fmt.Println(`read error:`, err)
		w.reader = nil
		return
	}

	w.fileSize += int64(n)

	read := w.buf[:n]

	w.Room.SendMsg(read)

	fmt.Println(`read`, n, `total`, w.fileSize)
}

func (w *world) Player(p rome.IPlayerConn, status bool) {

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
