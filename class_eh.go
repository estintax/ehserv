package main

import (
	"crypto/tls"
	"ehserv/dinolang"
	"fmt"
)

func EhClassHandler(args []string, segmentName string) bool {
	switch args[0] {
	case "default-mime", "default-page", "host", "webroot", "charset":
		if len(args) > 1 {
			if dinolang.GetTypeEx(args[1]) == "string" {
				if args[0] == "default-mime" {
					defaultMime = dinolang.StringToText(dinolang.IfVariableReplaceIt(args[1]).(string))
				} else if args[0] == "default-page" {
					if len(args) > 2 {
						if dinolang.GetTypeEx(args[2]) == "string" {
							vHostName := dinolang.StringToText(dinolang.IfVariableReplaceIt(args[2]).(string))
							if vHost, ok := vHosts[vHostName]; ok {
								vHost.DefaultPage = dinolang.StringToText(dinolang.IfVariableReplaceIt(args[1]).(string))
								vHosts[vHostName] = vHost
							} else {
								dinolang.PrintError("Unknown virtual host")
								dinolang.SetReturned("int", 0, segmentName)
								return false
							}
						} else {
							dinolang.PrintError("Second argument is not a string, is a " + dinolang.GetTypeEx(args[2]))
						}
					} else {
						defaultPage = dinolang.StringToText(dinolang.IfVariableReplaceIt(args[1]).(string))
					}
				} else if args[0] == "host" {
					ip = dinolang.StringToText(dinolang.IfVariableReplaceIt(args[1]).(string))
				} else if args[0] == "webroot" {
					webroot = dinolang.StringToText(dinolang.IfVariableReplaceIt(args[1]).(string))
				} else if args[0] == "charset" {
					charset = dinolang.StringToText(dinolang.IfVariableReplaceIt(args[1]).(string))
				}
				dinolang.SetReturned("int", 1, segmentName)
			} else {
				dinolang.PrintError("First argument is not a string, is a " + dinolang.GetTypeEx(args[1]))
				dinolang.SetReturned("int", 0, segmentName)
			}
		} else {
			dinolang.PrintError("Missed few arguments")
			dinolang.SetReturned("int", 0, segmentName)
		}
	case "port":
		if len(args) > 2 {
			if dinolang.GetTypeEx(args[1]) == "string" && dinolang.GetTypeEx(args[2]) == "int" {
				portType := dinolang.StringToText(dinolang.IfVariableReplaceIt(args[1]).(string))
				if portType == "http" {
					port = dinolang.IfVariableReplaceIt(args[2]).(int)
				} else if portType == "https" {
					tlsPort = dinolang.IfVariableReplaceIt(args[2]).(int)
				} else {
					dinolang.PrintError("Undefined port type")
					dinolang.SetReturned("int", 0, segmentName)
					return false
				}
				dinolang.SetReturned("int", 1, segmentName)
			} else {
				dinolang.PrintError("Some argument is not an int or a string")
				dinolang.SetReturned("int", 0, segmentName)
			}
		} else {
			dinolang.PrintError("Missed few arguments")
			dinolang.SetReturned("int", 0, segmentName)
		}
	case "add-vhost", "edit-vhost":
		if len(args) > 2 {
			if dinolang.GetTypeEx(args[1]) == "string" && dinolang.GetTypeEx(args[2]) == "string" {
				vHost := vHost{
					Root:        dinolang.StringToText(dinolang.IfVariableReplaceIt(args[2]).(string)),
					DefaultPage: defaultPage,
				}
				if len(args) > 4 {
					if dinolang.GetTypeEx(args[3]) == "string" && dinolang.GetTypeEx(args[4]) == "string" {
						cert, err := tls.LoadX509KeyPair(dinolang.StringToText(dinolang.IfVariableReplaceIt(args[3]).(string)), dinolang.StringToText(dinolang.IfVariableReplaceIt(args[4]).(string)))
						if err != nil {
							fmt.Printf("[WARN] [CONFIG] Failed to load certificate %s, skiping\n", dinolang.StringToText(dinolang.IfVariableReplaceIt(args[3]).(string)))
						} else {
							certs = append(certs, cert)
						}
					} else {
						dinolang.PrintError("Some argument is not a string")
					}
				}
				if !vHostsUsed {
					vHostsUsed = true
				}
				vHosts[dinolang.StringToText(dinolang.IfVariableReplaceIt(args[1]).(string))] = vHost
				dinolang.SetReturned("int", 1, segmentName)
			} else {
				dinolang.PrintError("Some argument is not a string")
				dinolang.SetReturned("int", 0, segmentName)
			}
		} else {
			dinolang.PrintError("Missed few arguments")
			dinolang.SetReturned("int", 0, segmentName)
		}
	}
	return true
}
