package main

import (
	"time"

	"github.com/zhengkai/rome"
)

var (
	worldPool  = make(map[string]*world)
	worldQuery = make(chan *worldCb, 1000)
)

type worldCb struct {
	id string
	ch chan rome.IWorld
}

func getWorld(id interface{}) rome.IWorld {

	ch := make(chan rome.IWorld)

	worldQuery <- &worldCb{
		id: id.(string),
		ch: ch,
	}

	w := <-ch

	return w
}

func worldDeamon() {

	tick := time.NewTicker(5 * time.Second)

	for {
		select {
		case t := <-tick.C:
			worldCheckIdle(t)
		case cb := <-worldQuery:
			w := worldCreate(cb.id)
			cb.ch <- w
		}
	}
}

func worldCreate(id string) (w *world) {

	w, ok := worldPool[id]
	if ok {
		return
	}

	w = &world{}
	w.init(id)

	worldPool[id] = w

	return
}

func worldCheckIdle(t time.Time) {

	t = t.Add(-15 * time.Second)

	for k, v := range worldPool {

		if v.online > 0 {
			continue
		}

		if v.timeEmpty.After(t) {
			continue
		}

		delete(worldPool, k)
		v.stop()
	}
}
