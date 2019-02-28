package main

import (
	"flag"
	"fmt"
	"net/http"
	"strconv"
)

func server() {

	port := flag.Int("port", 21002, "http port")
	flag.Parse()

	addr := `127.0.0.1:` + strconv.Itoa(*port)

	mux := http.NewServeMux()
	mux.HandleFunc(`/file`, doFile)
	mux.HandleFunc(`/listen`, doListen)
	fmt.Printf("port = %s\n", addr)
	err := http.ListenAndServe(addr, mux)
	if err != nil {
		fmt.Println(`http`, addr, `start fail:`)
		fmt.Println(err.Error())
	}
}
