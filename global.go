package main

import "crypto/tls"

const SERVER string = "EHServ/0.1"

var ip string = "0.0.0.0"
var port int = 8080
var charset string = "utf-8"
var webroot string = "/var/www"
var defaultMime string = "application/octet-stream"
var defaultPage string = "index.html"
var phpCgi string = "none"
var extEnum []string

var vHostsUsed bool = false
var vHosts map[string]string

var tlsPort int = 0
var certs []tls.Certificate
