package main

import (
	"bufio"
	"io"
	"log"
	"os"
	"time"

	"github.com/fsnotify/fsnotify"
)

type file struct {
	filename string
	ch       chan bool
	size     int64

	buf      []byte
	prevLine []byte
	reader   *bufio.Reader

	initOver bool
}

func (f *file) prev() (b []byte, ok bool) {

	if !f.initOver {
		return
	}

	size := int(f.size)

	reader, err := createReader(f.filename)
	if err != nil {
		return
	}

	load := 100000
	skip := 0

	if load > size {
		load = size
	} else {
		skip = int(size - load)
	}

	if skip > 0 {
		_, err := reader.Discard(skip)
		if err != nil {
			return
		}
	}

	b = make([]byte, load)

	n, err := reader.Read(b)
	if n != load {
		return
	}

	ok = true

	return
}

func (f *file) scan() (b []byte, reset bool, ok bool) {

	if !f.initOver {
		return
	}

	finfo, err := os.Stat(f.filename)
	if err != nil {
		f.reader = nil
		return
	}

	size := finfo.Size()

	if f.size > size {
		f.reader = nil
	}
	if f.reader == nil {
		f.reader, err = createReader(f.filename)
		if err != nil {
			return
		}
		f.size = 0
	}

	n, err := f.reader.Read(f.buf)
	if err != nil {
		if err != io.EOF {
			f.reader = nil
		}
		return
	}

	reset = f.size == 0
	f.size += int64(n)

	b = f.buf[:n]
	ok = true

	return
}

func (f *file) start() {

	go func() {

		f.buf = make([]byte, readBuffSize)

		finfo, err := os.Stat(f.filename)
		if err == nil {
			f.size = finfo.Size()
		}

		if f.size > 0 {
			f.reader, _ = createReader(f.filename)
			if f.reader != nil {
				f.reader.Discard(int(f.size))
			}
		}

		f.initOver = true
	}()

	go f.deamon()
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

	select {
	case f.ch <- true:
	default:
	}

	for {
		event := <-watcher.Events
		if event.Op&fsnotify.Write == fsnotify.Write {
			f.ch <- true
		}
	}
}
