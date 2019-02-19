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

func (p *player) LoginDecode(b []byte) (login interface{}, ok bool) {
	fmt.Println(`LoginDecode`)
	return nil, true
}

func (p *player) LoginAuth(login interface{}) (ok bool) {
	fmt.Println(`LoginAuth`)

	p.ID = 1

	return true
}

func (p *player) ParseMessage(b []byte) (msg interface{}, ok bool) {
	return string(b), true
}

func (p *player) GetID() interface{} {
	return p.ID
}
