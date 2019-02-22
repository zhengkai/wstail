package main

import (
	"fmt"
	"time"

	"github.com/zhengkai/rome"
)

var (
	buildTime      string
	buildGoVersion string

	oneWorld rome.IWorld
)

func main() {

	fmt.Println(`build time`, buildTime)
	// fmt.Println(`build by`, buildGoVersion)

	filePool.init()

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
	nw.init(id.(int))

	oneWorld = nw

	return nw
}
