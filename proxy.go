package main

import (
  "strings"
  "strconv"
  "time"
  "net"
  "fmt"
  "bufio"
  "crypto/tls"
)

func proxy(client net.Conn, headers string, url string) {
  protoAndUrl := strings.SplitN(url, "://", 2)
  hostAndUrl := strings.SplitN(protoAndUrl[1], "/", 2)
  var path string
  if len(hostAndUrl) == 1 {
    path = "/"
  } else {
    path = hostAndUrl[1]
  }

  var host string
  if strings.Contains(hostAndUrl[0], ":") == false {
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
    conn, err := net.DialTimeout("tcp", host, time.Second * 10)
    if err != nil {
      fmt.Println(err)
      return
    }

    handleProxyConnection(client, protoAndUrl[0], conn, host, headers, path)
  } else if protoAndUrl[0] == "https" {
    conn, err := tls.Dial("tcp", host, nil)
    if err != nil {
      fmt.Println(err)
      return
    }
    handleProxyConnection(client, protoAndUrl[0], conn, host, headers, path)
  }
}

func handleProxyConnection(client net.Conn, proto string, conn net.Conn, host string, clientRawHeaders string, url string) {
  var rawHeaders []string
  var splittedClientHeaders []string
  splittedClientHeaders = strings.Split(clientRawHeaders, "\r\n")
  if len(splittedClientHeaders) <= 1 {
    splittedClientHeaders = strings.Split(clientRawHeaders, "\n")
  }

  var originalHost string
  var handledClientHeaders string
  var findedBody int
  handledClientHeaders = splittedClientHeaders[0]
  for i := 0; i < len(splittedClientHeaders); i++ {
    if splittedClientHeaders[i] == "" {
      findedBody = i+1
      break
    }
    splitted := strings.SplitN(splittedClientHeaders[i], ": ", 2)
    if len(splitted) == 2 {
      if splitted[0] == "Host" {
        originalHost = splitted[1]
        splitted[1] = strings.SplitN(host, ":", 2)[0]
      }

      handledClientHeaders = handledClientHeaders + "\r\n" + splitted[0] + ": " + splitted[1]
    }
  }
  handledClientHeaders = handledClientHeaders + "\r\nX-Forwarded-For: " + strings.Split(client.RemoteAddr().String(), ":")[0] + "\r\n\r\n"
  if len(splittedClientHeaders) != findedBody {
    handledClientHeaders = handledClientHeaders + splittedClientHeaders[findedBody]
  }
  fmt.Fprintf(conn, handledClientHeaders)
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

  var contentLength int
  for i := 1; i < len(rawHeaders); i++ {
    splitted := strings.SplitN(rawHeaders[i], ": ", 2)
    if len(splitted) == 2 {
      if strings.ToLower(splitted[0]) == "Content-Length" {
        contentLength, _ = strconv.Atoi(splitted[1])
        break
      }
    }
  }

  for i := 0; i < len(rawHeaders); i++ {
    if strings.Contains(rawHeaders[i], "Host: ") {
      rawHeaders[i] = strings.Replace(rawHeaders[i], host, originalHost, 1)
    }
    fmt.Fprintf(client, "%s", rawHeaders[i])
  }
  fmt.Fprintf(client, "\r\n")

  body := make([]byte, contentLength)
  reader.Read(body)
  client.Write(body)

  conn.Close()
  client.Close()
}
