package main

import (
	"log"
	"time"

	"github.com/fsnotify/fsnotify"
)

type file struct {
	filename string
	ch       chan bool
}

func (f *file) deamon() {

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}

	f.ch = make(chan bool)

	for {
		err := watcher.Add(f.filename)
		if err == nil {
			break
		}
		time.Sleep(5 * time.Second)
	}

	for {
		event := <-watcher.Events
		if event.Op&fsnotify.Write == fsnotify.Write {
			f.ch <- true
		}
	}
}
