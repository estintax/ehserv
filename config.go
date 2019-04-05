package main

import (
	"fmt"
	"io"
	"os"
	"strings"
	"strconv"
)

func loadConfig(path string) bool {
	file, err := os.Open(path)
	if err != nil {
		fmt.Printf("[ERROR] Failed to open configuration file: %s\n", err.Error())
		return false
	}

	stat, _ := file.Stat()
	size := stat.Size()
	var conf string
	data := make([]byte, size)
	for {
		length, err := file.Read(data)
		if err == io.EOF {
			break
		}

		conf = string(data[:length])
	}

	file.Close()
	result := parseConfig(conf)
	return result
}

func parseConfig(conf string) bool {
	lines := strings.Split(conf, "\n")
	for i := 0; i < len(lines); i++ {
		params := strings.Split(lines[i], " ")
		if len(params) > 0 {
			switch params[0] {
			case "ip":
				ip = params[1]
			case "port":
				portInt, err := strconv.Atoi(params[1])
				if err != nil {
					fmt.Println("[WARN] [CONFIG] port is not a int, using default value")
				} else {
					port = portInt
				}
			case "charset":
				charset = params[1]
			case "webroot":
				webroot = params[1]
			case "default-type":
				defaultMime = params[1]
			case "default-indexpage":
				defaultPage = params[1]
			case "php-cgi":
				phpCgi = params[1]
			case "vhost":
				if len(params) >= 3 {
					if vHostsUsed == false {
						vHostsUsed = true
					}
					vHosts[params[1]] = params[2]
				} else {
					fmt.Println("[WARN] [CONFIG] vhost: missing parameters, skiping")
				}
			}
		}
	}

	return true
}