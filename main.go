package main

import (
	"fmt"
)

func main() {
	fmt.Println("EHServ 0.1 (c) 2019-2020, Maksim Pinigin")
	vHosts = make(map[string]string)
	if loadConfig("ehserv.conf") == true {
		startServer()
	}
}
