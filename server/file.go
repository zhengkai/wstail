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
	stop     bool
}

func (f *file) prev() (b []byte, ok bool) {

	if f.size == 0 {
		if f.initOver {
			return
		}
		time.Sleep(time.Second / 2)
		if f.size == 0 {
			return
		}
	}

	size := int(f.size)

	reader, err := createReader(f.filename)
	if err != nil {
		return
	}

	// 根据文件当前大小，确认之前应该发送多少
	// 如果 Discard 或者 ReadFull 失败，说明文件大小有变动（重置）
	// 则无须发送重置前的内容

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

	n, err := io.ReadFull(reader, b)
	if err != nil || n != load {
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

	f.ch = make(chan bool, 1)

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
	defer func() {
		watcher.Close()
		close(f.ch)
	}()

	for {
		err := watcher.Add(f.filename)
		if err == nil {
			break
		}
		if f.stop {
			return
		}
		time.Sleep(5 * time.Second)
	}

	tick := time.NewTicker(5 * time.Second)

	for {
		select {
		case <-tick.C:
			if f.stop {
				return
			}
		case event := <-watcher.Events:
			if event.Op&fsnotify.Write == fsnotify.Write {
				select {
				case f.ch <- true:
				default:
				}
			}
		}
	}
}
