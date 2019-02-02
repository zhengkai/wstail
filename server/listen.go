package main

import (
	"fmt"
	"net/http"
)

func doListen(w http.ResponseWriter, r *http.Request) {
	fmt.Println(`listen`)
}
