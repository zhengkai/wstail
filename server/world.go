package main

import (
	"fmt"

	"github.com/zhengkai/rome"
)

type world struct {
	rome.World
	id int
}

func (w *world) Tick(i int) (ok bool) {

	if i%gameFPS == 0 {
		w.Room.SendMsg([]byte(fmt.Sprintf(`server tick %d`, i)))
	}

	return true
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
