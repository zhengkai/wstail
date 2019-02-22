package main

import (
	"sync/atomic"

	"github.com/zhengkai/rome"
)

var (
	listenChan = make(chan rome.IPlayerConn, 1000)
	autoID     uint64
)

type player struct {
	rome.PlayerConn
}

func (p *player) Login(b []byte) (ok bool) {
	p.ID = atomic.AddUint64(&autoID, 1)
	return true
}

func (p *player) ParseMessage(b []byte) (msg interface{}, ok bool) {
	return string(b), true
}

func (p *player) GetID() interface{} {
	return p.ID
}
