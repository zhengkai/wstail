package main

import (
	"fmt"
	"sync"
)

var (
	filePool = &fileManager{
		pool: &sync.Map{},
	}
)

type fileManager struct {
	pool *sync.Map
	ch   chan *fileManagerCb
}

type fileManagerCb struct {
	filename string
	ch       chan *file
}

func (fm *fileManager) get(filename string) *file {

	v, ok := fm.pool.Load(filename)
	if ok {
		return v.(*file)
	}

	fmt.Println(`fm get 1`)

	cb := &fileManagerCb{
		filename: filename,
		ch:       make(chan *file),
	}

	fm.ch <- cb
	fmt.Println(`fm get 2`)

	f := <-cb.ch
	close(cb.ch)

	return f
}

func (fm *fileManager) init() {
	fm.ch = make(chan *fileManagerCb, 1000)
	go fm.deamon()
}

func (fm *fileManager) deamon() {

	for {
		fmt.Println(`fm get 4`)
		cb := <-fm.ch
		fmt.Println(`fm get 3`)

		v, ok := fm.pool.Load(cb.filename)
		if ok {
			cb.ch <- v.(*file)
		}

		f := &file{
			filename: cb.filename,
		}
		f.start()

		fm.pool.Store(cb.filename, f)

		cb.ch <- f
	}
}
