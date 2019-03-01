package main

import (
	"flag"
	"fmt"
	"net/http"
)

func server() {

	flag.Parse()

	addr := fmt.Sprintf(`%s:%d`, listenHost, *listenPort)

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
