package main

import (
	"crypto/tls"
	"fmt"
	"strconv"
)

func startHTTPSServer() {
	cfg := &tls.Config{Certificates: certs}

	portStr := strconv.Itoa(tlsPort)
	addr := ip + ":" + portStr
	tlsSvr, err := tls.Listen("tcp4", addr, cfg)
	if err != nil {
		fmt.Printf("[ERROR] [HTTPS] Failed to start HTTPS server: %s\n", err.Error())
		return
	}

	for {
		conn, err := tlsSvr.Accept()
		if err != nil {
			fmt.Printf("[ERROR] [HTTPS] Failed to accept connection: %s\n", err.Error())
			continue
		}

		go connectionHandler(conn, true)
	}
}
