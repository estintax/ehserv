package main

import (
	"fmt"
	"net"
	"strconv"
	"bufio"
	"strings"
	"os"
	"os/exec"
	"io"
)

var isPHP bool = false

func startServer(ip string, port int) {
	portStr := strconv.Itoa(port)
	addr := ip + ":" + portStr
	serv, err := net.Listen("tcp4", addr)
	if err != nil {
		fmt.Printf("[ERROR] [SERVER] Failed to start server: %s\n", err.Error())
		return
	}

	for {
		conn, err := serv.Accept()
		if err != nil {
			fmt.Printf("[ERROR] [SERVER] Failed to accept connection: %s\n", err.Error())
			return
		}

		go connectionHandler(conn)
	}
}

func connectionHandler(conn net.Conn) {
	var buffer string

	reader := bufio.NewReader(conn)
	for {
		str, err := reader.ReadString('\n')
		if err != nil {
			fmt.Fprintf(conn, "HTTP/1.1 500 Internal Server Error\nServer: %s\n\n500 Internal Server Error\n", SERVER);
			return
		}

		buffer = buffer + str
		if str == "\r\n" || str == "\n" {
			lines := strings.Split(buffer, "\n")
			for i := 0; i < len(lines); i++ {
				if i == 0 {
					if lines[i] == "" || lines[i] == "\r" {
						fmt.Fprintf(conn, "HTTP/1.1 404 Bad Request\nServer: %s\n\n<h1>400 Bad Request</h1>\n", SERVER);
						conn.Close()
						return
					}
					params := strings.Split(lines[i], " ")
					fmt.Printf("[INFO] [%s] %s - %s\n", conn.RemoteAddr().String(), params[0], params[1])

					handleURL(conn, params[1])
					break
				}
			}
		}
	}
}

func sendHTTPResponse(conn net.Conn, code int, contentType string, content interface{}) {
	var codeStr string
	var contentLength int

	switch code {
	case 200: codeStr = "HTTP/1.1 200 OK"
	case 301: codeStr = "HTTP/1.1 301 Moved Permanently"
	case 302: codeStr = "HTTP/1.1 302 Found"
	case 400: codeStr = "HTTP/1.1 400 Bad Request"
	case 403: codeStr = "HTTP/1.1 403 Forbidden"
	case 404: codeStr = "HTTP/1.1 404 Not Found"
	case 500: codeStr = "HTTP/1.1 500 Internal Server Error"
	default: codeStr = "HTTP/1.1 418 I'm a teapot"
	}

	switch content.(type) {
	case []byte:
		contentLength = len(content.([]byte))
	case string:
		contentLength = len(content.(string))
	default:
		contentLength = len(content.(string))
	}

	if isPHP == true {
		fmt.Fprintf(conn, "%s\nServer: %s\nConnection: close\n", codeStr, SERVER)
	} else {
		fmt.Fprintf(conn, "%s\nContent-Type: %s; charset=%s\nContent-Length: %d\nServer: %s\nConnection: close\n\n", codeStr, contentType, charset, contentLength, SERVER)
	}

	switch content.(type) {
	case []byte:
		conn.Write(content.([]byte))
	case string:
		fmt.Fprint(conn, content.(string))
	default:
		fmt.Fprint(conn, content.(string))
	}

	isPHP = false
	conn.Close()
}

func handleURL(conn net.Conn, urlp string) bool {
	var url string
	//var origUrl string
	if urlp[len(urlp)-1] == '/' {
		url = urlp + defaultPage
	} else {
		url = urlp
	}

	if questionIndex := strings.Index(url, "?"); questionIndex != -1 {
		//origUrl = url
		url = url[:questionIndex]
	}

	file, err := os.Open(webroot + url)
	if err != nil {
		sendHTTPResponse(conn, 404, "text/html", "<h1>404 Not Found</h1>")
		return false
	}

	stat, _ := file.Stat()

	if stat.IsDir() == true {
		sendHTTPResponse(conn, 403, "text/html", "<h1>403 Forbidden</h1>")
		file.Close()
		return false
	}

	data := make([]byte, stat.Size())
	for {
		_, err := file.Read(data)
		if err == io.EOF {
			break
		}
	}

	file.Close()

	parseOne := strings.Split(url, "/")
	parseTwo := strings.Split(parseOne[len(parseOne)-1], ".")
	format := parseTwo[len(parseTwo)-1]

	if format == "php" && phpCgi != "none" {
		cmd := exec.Command(phpCgi)
		stdin, _ := cmd.StdinPipe()
		stdout, _ := cmd.StdoutPipe()
		err = cmd.Start()
		if err != nil {
			sendHTTPResponse(conn, 500, "text/html", "<h1>500 Internal Server Error</h1></h2>Failed to start PHP Server: " + err.Error() + "</h2>")
			return false
		}
		stdin.Write(data)
		phpReader := bufio.NewReader(stdout)
		//length := reader.Size()
		//phpData := make([]byte, length)
		//length, _ = reader.Read(phpData)
		var phpData string
		for {
			singleByte, err := phpReader.ReadByte()
			if err != nil {
				break
			}

			phpData += string(singleByte)
		}
		isPHP = true
		//sendHTTPResponse(conn, 200, "", string(phpData[:length]))
		sendHTTPResponse(conn, 200, "", phpData)
		phpData = ""
		data = nil
		return true
	}
	
	sendHTTPResponse(conn, 200, getContentType(format), data)
	data = nil
	return true
}