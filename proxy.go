package main

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"io"
	"net"
	"strings"
	"time"
)

func proxy(client net.Conn, headers string, url string, proxyId int, isTLS bool) bool {
	protoAndUrl := strings.SplitN(url, "://", 2)
	hostAndUrl := strings.SplitN(protoAndUrl[1], "/", 2)
	var path string
	if len(hostAndUrl) == 1 {
		path = "/"
	} else {
		path = hostAndUrl[1]
	}

	var host string
	if !strings.Contains(hostAndUrl[0], ":") {
		switch protoAndUrl[0] {
		case "http":
			host = hostAndUrl[0] + ":80"
		case "https":
			host = hostAndUrl[0] + ":443"
		}
	} else {
		host = hostAndUrl[0]
	}

	if protoAndUrl[0] == "http" {
		conn, err := net.DialTimeout("tcp", host, time.Second*10)
		if err != nil {
			fmt.Println(err)
			return false
		}

		handleProxyConnection(client, protoAndUrl[0], conn, host, headers, path, proxyId, isTLS)
	} else if protoAndUrl[0] == "https" {
		conn, err := tls.Dial("tcp", host, nil)
		if err != nil {
			fmt.Println(err)
			return false
		}
		handleProxyConnection(client, protoAndUrl[0], conn, host, headers, path, proxyId, isTLS)
	}

	return true
}

func handleProxyConnection(client net.Conn, proto string, conn net.Conn, host string, clientRawHeaders string, url string, proxyId int, isTLS bool) {
	var rawHeaders []string
	var splittedClientHeaders []string
	splittedClientHeaders = strings.Split(clientRawHeaders, "\r\n")
	if len(splittedClientHeaders) <= 1 {
		splittedClientHeaders = strings.Split(clientRawHeaders, "\n")
	}

	var originalHost string
	var handledClientHeaders string
	var findedBody int
	splittedMainHead := strings.Split(splittedClientHeaders[0], " ")
	targetPath := "/" + strings.SplitN(proxyUrls[proxyId].Address, "/", 4)[3]
	splittedMainHead[1] = strings.Replace(splittedMainHead[1], splittedMainHead[1][0:len(proxyUrls[proxyId].Url)], targetPath, 1)
	handledClientHeaders = fmt.Sprintf("%s %s %s", splittedMainHead[0], splittedMainHead[1], splittedMainHead[2])
	for i := 0; i < len(splittedClientHeaders); i++ {
		if splittedClientHeaders[i] == "" {
			findedBody = i + 1
			break
		}
		splitted := strings.SplitN(splittedClientHeaders[i], ": ", 2)
		if len(splitted) == 2 {
			if splitted[0] == "Host" {
				originalHost = splitted[1]
				splitted[1] = strings.SplitN(host, ":", 2)[0]
			} else if splitted[0] == "Connection" {
				splitted[1] = "close"
			}

			handledClientHeaders = handledClientHeaders + "\r\n" + splitted[0] + ": " + splitted[1]
		}
	}
	handledClientHeaders = handledClientHeaders + "\r\nX-Forwarded-For: " + strings.Split(client.RemoteAddr().String(), ":")[0] + "\r\n\r\n"
	if len(splittedClientHeaders) != findedBody {
		handledClientHeaders = handledClientHeaders + splittedClientHeaders[findedBody]
	}
	fmt.Fprint(conn, handledClientHeaders)
	reader := bufio.NewReader(conn)
	for {
		readed, err := reader.ReadString('\n')
		if err != nil {
			break
		}

		if readed == "\r\n" || readed == "\n" {
			break
		}

		rawHeaders = append(rawHeaders, readed)
	}

	for i := 0; i < len(rawHeaders); i++ {
		splittedHeader := strings.Split(rawHeaders[i], ": ")
		if strings.ToLower(splittedHeader[0]) == "server" {
			rawHeaders[i] = splittedHeader[0] + ": " + SERVER + "\r\n"
		} else if strings.ToLower(splittedHeader[0]) == "location" {
			splitted := strings.Split(splittedHeader[1], "/")
			protoStr := "http://"
			if isTLS {
				protoStr = "https://"
			}
			rawHeaders[i] = strings.Replace(rawHeaders[i], fmt.Sprintf("%s//%s", splitted[0], splitted[2]), fmt.Sprintf("%s%s%s", protoStr, originalHost, proxyUrls[proxyId].Url), 1)
		} else {
			rawHeaders[i] = strings.ReplaceAll(rawHeaders[i], strings.SplitN(host, ":", 2)[0], originalHost)
			if isTLS {
				rawHeaders[i] = strings.ReplaceAll(rawHeaders[i], "http://", "https://")
			} else {
				rawHeaders[i] = strings.ReplaceAll(rawHeaders[i], "https://", "http://")
			}
		}
		fmt.Fprintf(client, "%s", rawHeaders[i])
	}
	fmt.Fprintf(client, "\r\n")

	go func() {
		_, err := reader.WriteTo(client)
		if err != nil {
			reader.Reset(client)
		}
		client.Close()
	}()

	go func() {
		io.Copy(conn, client)
		conn.Close()
	}()
}
