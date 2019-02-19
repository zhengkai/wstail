package main

import "fmt"

type conn struct {
	id   int
	name string
}

func (c *conn) add() {
	c.id++
}

func (c *conn) get() int {
	return c.id
}

type abc struct {
	conn
}

func (c *abc) get() int {
	return c.id + 5
}

func test() {

	c := &abc{}
	c.id = 100

	c.add()
	fmt.Println(c.get())
	fmt.Println(c.conn.get())
}
