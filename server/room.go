package main

import (
	"time"

	"github.com/zhengkai/rome"
)

type room struct {
	rome.Room
}

func (r *room) GetTickDuration() time.Duration {
	return time.Second / time.Duration(gameFPS)
}
