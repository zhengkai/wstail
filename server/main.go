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

	// getWorld(1)

	go worldDeamon()

	go rome.Manager(listenChan, getWorld)

	go server()
	time.Sleep(time.Hour * 999999)
}
