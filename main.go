package main

import (
	"fmt"
)

func main() {
	fmt.Println("EHServ 0.1 (c) 2019, Maksim Pinigin")
	if loadConfig("ehserv.conf") == true {
		startServer(ip, port)
	}
}