package main

import (
	"fmt"
)

func main() {
	fmt.Println("EHServ 0.1 (c) 2019-2021, Maksim Pinigin")
	fillStatusCodes()
	vHosts = make(map[string]string)
	if loadConfig("ehserv.conf") {
		startServer()
	}
}
