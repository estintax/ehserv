package main

import (
	"fmt"
	"io"
	"os"
	"strings"
	"strconv"
	"crypto/tls"
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
				if len(params) >= 2 {
					ip = params[1]
				} else {
					fmt.Println("[WARN] [CONFIG] ip: missing parameters, skiping")
				}
			case "port":
				if len(params) >= 2 {
					portInt, err := strconv.Atoi(params[1])
					if err != nil {
						fmt.Println("[WARN] [CONFIG] port is not a int, using default value")
					} else {
						port = portInt
					}
				} else {
					fmt.Println("[WARN] [CONFIG] port: missing parameters, skiping")
				}
			case "charset":
				if len(params) >= 2 {
					charset = params[1]
				} else {
					fmt.Println("[WARN] [CONFIG] charset: missing parameters, skiping")
				}
			case "webroot":
				if len(params) >= 2 {
					webroot = params[1]
				} else {
					fmt.Println("[WARN] [CONFIG] webroot: missing parameters, skiping")
				}
			case "default-type":
				if len(params) >= 2 {
					defaultMime = params[1]
				} else {
					fmt.Println("[WARN] [CONFIG] default-type: missing parameters, skiping")
				}
			case "default-indexpage":
				if len(params) >= 2 {
					defaultPage = params[1]
				} else {
					fmt.Println("[WARN] [CONFIG] default-indexpage: missing parameters, skiping")
				}
			case "cgi":
				if len(params) >= 3 {
					extEnum = strings.Split(params[2], ",")
					phpCgi = params[1]
				} else {
					fmt.Println("[WARN] [CONFIG] cgi: missing parameters, skiping")
				}
			case "vhost":
				if len(params) >= 3 {
					if vHostsUsed == false {
						vHostsUsed = true
					}
					vHosts[params[1]] = params[2]
					if len(params) >= 5 {
						cert, err := tls.LoadX509KeyPair(params[3], params[4])
						if err != nil {
							fmt.Printf("[WARN] [CONFIG] Failed to load certificate %s, skiping\n", params[3])
							continue
						}
						certs = append(certs, cert)
					}
				} else {
					fmt.Println("[WARN] [CONFIG] vhost: missing parameters, skiping")
				}
			case "ssl":
				if len(params) >= 2 {
					portInt, err := strconv.Atoi(params[1])
					if err != nil {
						fmt.Println("[WARN] [CONFIG] ssl port is not a int, ssl disabled")
					} else {
						tlsPort = portInt
					}
				} else {
					fmt.Println("[WARN] [CONFIG] ssl: missing parameters, skiping")
				}
			case "proxy":
				if len(params) >= 3 {
					proxyUrls = append(proxyUrls, Proxy{Url: params[2], Vhost: params[1], Address: params[3]})
				} else {
					fmt.Println("[WARN] [CONFIG] proxy: missing parameters, skiping")
				}
			}
		}
	}

	return true
}
