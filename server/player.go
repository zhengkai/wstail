package main

import (
	"fmt"

	"github.com/zhengkai/rome"
)

var (
	listenChan = make(chan rome.IPlayerConn, 1000)
)

type player struct {
	rome.PlayerConn
}

func (p *player) Login(b []byte) (ok bool) {
	fmt.Println(`Login v2`)
	p.ID = 1
	return true
}

func (p *player) ParseMessage(b []byte) (msg interface{}, ok bool) {
	return string(b), true
}

func (p *player) GetID() interface{} {
	return p.ID
}
