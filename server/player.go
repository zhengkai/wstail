package main

import (
	"sync/atomic"

	"pb"

	"github.com/gogo/protobuf/proto"
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

	login := &pb.Login{}

	err := proto.Unmarshal(b, login)
	if err != nil {
		return
	}

	if !checkFileName(login.FileName) {
		return
	}

	p.ID = atomic.AddUint64(&autoID, 1)
	p.file = dirBase + `/` + login.FileName

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
