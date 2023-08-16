package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func startServer() {
	if tlsPort != 0 {
		go startHTTPSServer()
	}

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
			continue
		}

		go connectionHandler(conn, false)
	}
}

func connectionHandler(conn net.Conn, isTLS bool) {

	var buffer string = ""
	var firstLine string = ""
	var contentlength int
	var query string

	reader := bufio.NewReader(conn)
	for {
		str, err := reader.ReadString('\n')
		if err != nil {
			sendHTTPResponse(conn, 500, "text/plain", "500 Internal Server Error", false)
			return
		}

		if firstLine == "" && buffer == "" {
			firstLine = strings.Split(str, " ")[0]
		}

		buffer = buffer + str
		if str == "\r\n" && firstLine != "POST" || str == "\n" && firstLine != "POST" {
			lines := strings.Split(buffer, "\r\n")
			if len(lines) <= 1 {
				lines = strings.Split(buffer, "\n")
			}
			for i := 0; i < len(lines); i++ {
				if i == 0 {
					if lines[i] == "" || lines[i] == "\n" || lines[i] == "\r\n" {
						sendHTTPResponse(conn, 400, "text/html", "<h1>400 Bad Request</h1>", false)
						return
					}
					params := strings.Split(lines[i], " ")
					fmt.Printf("[INFO] [%s] %s - %s\n", conn.RemoteAddr().String(), params[0], params[1])

					handleURL(conn, params[0], params[1], lines, "", isTLS)
					break
				}
			}
		} else if firstLine == "POST" {
			if strings.Contains(str, ": ") {
				param := strings.SplitN(str, ": ", 2)
				if param[0] == "Content-Length" {
					contentlength, _ = strconv.Atoi(strings.Replace(param[1], "\r\n", "", -1))
				}
			}
			if str == "\r\n" {
				bytes := make([]byte, contentlength)
				reader.Read(bytes)
				query = string(bytes[:contentlength])
				params := strings.Split(strings.Split(buffer, "\r\n")[0], " ")
				lines := strings.Split(buffer, "\r\n")
				fmt.Printf("[INFO] [%s] %s - %s\n", conn.RemoteAddr().String(), params[0], params[1])
				handleURL(conn, params[0], params[1], lines, query, isTLS)
				break
			}
		}
	}
}

func sendHTTPResponse(conn net.Conn, code int, contentType string, content interface{}, isPHP bool) {
	var codeStr string
	var contentLength int

	codeStr = "HTTP/1.1 " + statusCodes[code]

	switch content.(type) {
	case []byte:
		contentLength = len(content.([]byte))
	case string:
		contentLength = len(content.(string))
	default:
		contentLength = len(content.(string))
	}

	if isPHP {
		fmt.Fprintf(conn, "%s\r\n", codeStr)
	} else {
		fmt.Fprintf(conn, "%s\r\nDate: %s\r\nCache-Control: no-cache\r\nContent-Type: %s; charset=%s\r\nContent-Length: %d\r\nServer: %s\r\nConnection: close\r\n\r\n", codeStr, getDate(), contentType, charset, contentLength, SERVER)
	}

	switch content.(type) {
	case []byte:
		conn.Write(content.([]byte))
	case string:
		fmt.Fprint(conn, content.(string))
	default:
		fmt.Fprint(conn, content.(string))
	}

	conn.Close()
}

func handleURL(conn net.Conn, method string, urlp string, all []string, query string, isTLS bool) bool {
	var url string
	var urlWithQuestion string
	var host string = ""
	var docroot string
	urlWithQuestion = ""
	docroot = webroot
	//var origUrl string
	if urlp[len(urlp)-1] == '/' {
		url = urlp + defaultPage
	} else {
		url = urlp
	}

	if questionIndex := strings.Index(url, "?"); questionIndex != -1 {
		//origUrl = url
		urlWithQuestion = url
		url = url[:questionIndex]
	}

	for i := 0; i < len(all); i++ {
		if strings.Contains(all[i], ": ") {
			param := strings.SplitN(all[i], ": ", 2)
			if param[0] == "Host" {
				host = param[1]
			}
		}
	}

	if host == "" {
		sendHTTPResponse(conn, 400, "text/html", "<h1>400 Bad Request</h1>", false)
		return false
	}

	if len(proxyUrls) != 0 {
		for i := 0; i < len(proxyUrls); i++ {
			if proxyUrls[i].Vhost == host {
				if strings.Contains(url, proxyUrls[i].Url) {
					if !proxy(conn, strings.Join(all, "\r\n")+"\r\n"+query, proxyUrls[i].Address, i, isTLS) {
						sendHTTPResponse(conn, 503, "text/html", "<h1>503 Bad Gateway<h1>", false)
					}
					return true
				}
			}
		}
	}

	if vHostsUsed && vHosts[host] != "" {
		docroot = vHosts[host]
	}

	file, err := os.Open(docroot + url)
	if err != nil {
		sendHTTPResponse(conn, 404, "text/html", "<h1>404 Not Found</h1>", false)
		return false
	}
	defer file.Close()

	stat, _ := file.Stat()

	if stat.IsDir() {
		sendHTTPResponse(conn, 403, "text/html", "<h1>403 Forbidden</h1>", false)
		file.Close()
		return false
	}

	parseOne := strings.Split(url, "/")
	parseTwo := strings.Split(parseOne[len(parseOne)-1], ".")
	format := parseTwo[len(parseTwo)-1]

	if findExt(format) && phpCgi != "none" {
		//var query string
		params := make(map[string]string)
		for i := 0; i < len(all); i++ {
			if i == 0 {
				continue
			} else if all[i] == "" || all[i] == "\r\n" {
				continue
			} else if i == len(all)-1 {
				break
			}

			splitted := strings.SplitN(all[i], ": ", 2)
			params[splitted[0]] = splitted[1]
		}

		data := make([]byte, stat.Size())
		for {
			_, err := file.Read(data)
			if err == io.EOF {
				break
			}
		}
		file.Close()

		cmd := exec.Command(phpCgi, docroot+url)
		cmd.Env = append(os.Environ(),
			"REMOTE_ADDR="+strings.Split(conn.RemoteAddr().String(), ":")[0],
			"REMOTE_HOST="+strings.Split(conn.RemoteAddr().String(), ":")[0],
			"REQUEST_METHOD="+method,
			"SERVER_NAME="+ip,
			"SERVER_PORT="+strconv.Itoa(port),
			"SERVER_PROTOCOL=HTTP/1.1",
			"SERVER_SOFTWARE="+SERVER,
			"CONTENT_LENGTH="+params["Content-Length"],
			"CONTENT_TYPE="+params["Content-Type"],
			"REDIRECT_STATUS=1",
			"SCRIPT_FILENAME="+docroot+url)

		if urlWithQuestion != "" {
			cmd.Env = append(cmd.Env, "QUERY_STRING="+strings.SplitN(urlWithQuestion, "?", 2)[1])
		}

		for i := 0; i < len(all); i++ {
			if strings.Contains(all[i], ": ") {
				param := strings.SplitN(all[i], ": ", 2)
				cmd.Env = append(cmd.Env, strings.ToUpper("HTTP_"+param[0]+"=")+param[1])
			}
		}

		stdin, _ := cmd.StdinPipe()
		if query != "" {
			stdin.Write([]byte(query))
		}
		stdout, _ := cmd.StdoutPipe()
		err = cmd.Start()
		if err != nil {
			sendHTTPResponse(conn, 502, "text/html", "<h1>502 Bad Gateway</h1>", false)
			return false
		}
		stdin.Write(data)
		phpReader := bufio.NewReader(stdout)
		var phpData string
		for {
			singleByte, err := phpReader.ReadByte()
			if err != nil {
				break
			}

			phpData += string(singleByte)
		}
		err = cmd.Wait()
		if err != nil {
			sendHTTPResponse(conn, 502, "text/html", "<h1>502 Bad Gateway</h1>", false)
			return true
		}
		var statusCode int = 200
		splittedHeaders := strings.Split(string(phpData), "\r\n")
		if len(splittedHeaders) <= 1 {
			splittedHeaders = strings.Split(string(phpData), "\n")
		}
		for i := 0; i < len(splittedHeaders); i++ {
			if strings.Contains(splittedHeaders[i], "Status: ") {
				statusCode, _ = strconv.Atoi(strings.SplitN(splittedHeaders[i], " ", 3)[1])
			} else if splittedHeaders[i] == "" {
				break
			}
		}
		headersSplit := strings.SplitN(phpData, "\r\n\r\n", 2)
		if len(headersSplit) == 1 {
			headersSplit = strings.SplitN(phpData, "\n\n", 2)
		}
		if !strings.Contains(headersSplit[0], "Server: ") {
			phpData = "Server: " + SERVER + "\r\n" + phpData
		}
		if !strings.Contains(headersSplit[0], "Connection: ") {
			phpData = "Connection: close\r\n" + phpData
		}
		if !strings.Contains(headersSplit[0], "Cache-Control: ") {
			phpData = "Cache-Control: no-cache\r\n" + phpData
		}
		if !strings.Contains(headersSplit[0], "Date: ") {
			phpData = "Date: " + getDate() + "\r\n" + phpData
		}
		sendHTTPResponse(conn, statusCode, "", phpData, true)
		return true
	}

	contentType := getContentType(format)
	if contentType == "text/plain" || contentType == "text/html" {
		contentType = contentType + "; charset=" + charset
	}
	fmt.Fprintf(conn, "HTTP/1.1 200 OK\r\nDate: %s\r\nCache-Control: no-cache\r\nContent-Type: %s\r\nContent-Length: %d\r\nServer: %s\r\nConnection: close\r\n\r\n", getDate(), contentType, stat.Size(), SERVER)
	defer conn.Close()
	fileReader := bufio.NewReader(file)
	fileReader.WriteTo(conn)
	return true
}
