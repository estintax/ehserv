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

		go connectionHandler(conn)
	}
}

func connectionHandler(conn net.Conn) {

	var buffer string = ""
	var firstLine string = ""
	var contentlength int
	var query string

	reader := bufio.NewReader(conn)
	for {
		str, err := reader.ReadString('\n')
		if err != nil {
			fmt.Fprintf(conn, "HTTP/1.1 500 Internal Server Error\nServer: %s\n\n500 Internal Server Error\n", SERVER);
			return
		}

		if firstLine == "" && buffer == "" {
			firstLine = strings.Split(str, " ")[0]
		}

		buffer = buffer + str
		if str == "\r\n" && firstLine != "POST" || str == "\n" && firstLine != "POST" {
			strings.Trim(buffer, "\r")
			lines := strings.Split(buffer, "\r\n")
			for i := 0; i < len(lines); i++ {
				if i == 0 {
					if lines[i] == "" || lines[i] == "\r" {
						fmt.Fprintf(conn, "HTTP/1.1 404 Bad Request\nServer: %s\n\n<h1>400 Bad Request</h1>\n", SERVER);
						conn.Close()
						return
					}
					params := strings.Split(lines[i], " ")
					fmt.Printf("[INFO] [%s] %s - %s\n", conn.RemoteAddr().String(), params[0], params[1])

					handleURL(conn, params[0], params[1], lines, "")
					break
				}
			}
		} else if firstLine == "POST" {
			if strings.Index(str, ": ") != -1 {
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
				handleURL(conn, params[0], params[1], lines, query)
				break
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

func handleURL(conn net.Conn, method string, urlp string, all []string, query string) bool {
	var url string
	var urlWithQuestion string
	var host string
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

	for i := 0; i<len(all); i++ {
		if strings.Index(all[i], ": ") != -1 {
			param := strings.SplitN(all[i], ": ", 2)
			if param[0] == "Host" {
				host = param[1]
			}
		}
	}

	if vHostsUsed == true && vHosts[host] != "" {
		docroot = vHosts[host]
	}

	file, err := os.Open(docroot + url)
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

	if findExt(format) == true && phpCgi != "none" {
		var params map[string]string
		//var query string
		params = make(map[string]string)
		for i := 0; i<len(all); i++ {
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
		cmd := exec.Command(phpCgi)
		cmd.Env = append(os.Environ(),
			"REMOTE_ADDR=" + strings.Split(conn.RemoteAddr().String(), ":")[0],
			"REMOTE_HOST=" + strings.Split(conn.RemoteAddr().String(), ":")[0],
			"REQUEST_METHOD=" + method,
			"SERVER_NAME=" + ip,
			"SERVER_PORT=" + strconv.Itoa(port),
			"SERVER_PROTOCOL=HTTP/1.1",
			"SERVER_SOFTWARE=" + SERVER,
			"CONTENT_LENGTH=" + params["Content-Length"],
			"CONTENT_TYPE=" + params["Content-Type"],
			"REDIRECT_STATUS=1",
			"SCRIPT_FILENAME=" + docroot + url)

		if urlWithQuestion != "" {
			cmd.Env = append(cmd.Env, "QUERY_STRING=" + strings.SplitN(urlWithQuestion, "?", 2)[1])
		}
		if query != "" {
			cmd.Env = append(cmd.Env, "QUERY_STRING=" + query)
		}

		for i := 0; i<len(all); i++ {
			if strings.Index(all[i], ": ") != -1 {
				param := strings.SplitN(all[i], ": ", 2)
				cmd.Env = append(cmd.Env, strings.ToUpper("HTTP_" + param[0] + "=") + param[1])
			}
		}

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
