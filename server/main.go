package main

import (
	"time"

	"github.com/zhengkai/rome"
)

var (
	oneWorld rome.IWorld
)

func main() {

	getWorld(1)

	go rome.Manager(listenChan, getWorld)

	go server()
	time.Sleep(time.Hour * 999999)
}

func getWorld(id interface{}) rome.IWorld {

	if oneWorld != nil {
		return oneWorld
	}

	nw := &world{}
	nw.id = id.(int)

	r := &room{}
	r.World = nw

	nw.Room = r

	rome.InitRoom(r)

	oneWorld = nw

	return nw
}
