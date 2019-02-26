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
	file string
}

func (p *player) Login(b []byte) (ok bool) {
	p.ID = atomic.AddUint64(&autoID, 1)
	p.file = `/tmp/fortune.txt`
	return true
}

func (p *player) GetWorldID() interface{} {
	return p.file
}

func (p *player) ParseMessage(b []byte) (msg interface{}, ok bool) {
	return string(b), true
}

func (p *player) GetID() interface{} {
	return p.ID
}
