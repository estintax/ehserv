package main

import (
	"ehserv/dinolang"
	"fmt"
	"os"
)

func main() {
	fmt.Println("EHServ 0.1 (c) 2019-2021, Maksim Pinigin")
	fillStatusCodes()
	vHosts = make(map[string]vHost)
	dinolang.Classes["eh"] = dinolang.Class{
		Prefix:    "eh",
		Used:      true,
		IsBuiltIn: false,
		Caller:    EhClassHandler,
		Loader:    nil}
	if len(os.Args) > 1 && os.Args[1] == "--config-cli" {
		dinolang.PiniginShell()
	}
	if dinolang.ParseFile("ehserv.conf.dino") {
		fmt.Println("Configuration script completed")
	}
	if loadConfig("ehserv.conf") {
		startServer()
	}
}
