package main

var statusCodes = make(map[int]string)

func fillStatusCodes() {
	// 100
	statusCodes[100] = "100 Continue"
	statusCodes[101] = "101 Switching Protocols"
	statusCodes[102] = "102 Processing"
	statusCodes[103] = "103 Early Hints"
	// 200
	statusCodes[200] = "200 OK"
	statusCodes[201] = "201 Created"
	statusCodes[202] = "202 Accepted"
	statusCodes[203] = "203 Non-Authoritative Information"
	statusCodes[204] = "204 No Content"
	statusCodes[205] = "205 Reset Content"
	statusCodes[206] = "206 Partial Content"
	statusCodes[207] = "207 Multi-Status"
	statusCodes[208] = "208 Already Reported"
	statusCodes[209] = "209 IM Used"
	// 300
	statusCodes[300] = "300 Multiple Choices"
	statusCodes[301] = "301 Moved Permanently"
	statusCodes[302] = "302 Found"
	statusCodes[303] = "303 See Other"
	statusCodes[304] = "304 Not Modified"
	statusCodes[305] = "305 Use Proxy"
	statusCodes[306] = "306 Switch Proxy"
	statusCodes[307] = "307 Temporary Redirect"
	statusCodes[308] = "308 Permanent Redirect"
	// 400
	statusCodes[400] = "400 Bad Request"
	statusCodes[401] = "401 Unauthorized"
	statusCodes[402] = "402 Payment Required"
	statusCodes[403] = "403 Forbidden"
	statusCodes[404] = "404 Not Found"
	statusCodes[405] = "405 Method Not Allowed"
	statusCodes[406] = "406 Not Acceptable"
	statusCodes[407] = "407 Proxy Authentication Required"
	statusCodes[408] = "408 Request Timeout"
	statusCodes[409] = "409 Conflict"
	statusCodes[410] = "410 Gone"
	statusCodes[411] = "411 Length Required"
	statusCodes[412] = "412 Precondition Failed"
	statusCodes[413] = "413 Payload Too Large"
	statusCodes[414] = "414 URI Too Long"
	statusCodes[415] = "415 Unsupported Media Type"
	statusCodes[416] = "416 Range Not Satisfiable"
	statusCodes[417] = "417 Expectation Failed"
	statusCodes[418] = "418 I'm a teapot"
	statusCodes[421] = "421 Misdirected Request"
	statusCodes[422] = "422 Unprocessable Entity"
	statusCodes[423] = "423 Locked"
	statusCodes[424] = "424 Failed Dependency"
	statusCodes[425] = "425 Too Early"
	statusCodes[426] = "426 Upgrade Required"
	statusCodes[428] = "428 Precondition Required"
	statusCodes[429] = "429 Too Many Requests"
	statusCodes[431] = "431 Request Header Fields Too Large"
	statusCodes[449] = "449 Retry With"
	statusCodes[451] = "451 Unavailable For Legal Reasons"
	statusCodes[499] = "499 Client Closed Request"
	// 500
	statusCodes[500] = "500 Internal Server Error"
	statusCodes[501] = "501 Not Implemented"
	statusCodes[502] = "502 Bad Gateway"
	statusCodes[503] = "503 Service Unavailable"
	statusCodes[504] = "504 Gateway Timeout"
	statusCodes[505] = "505 HTTP Version Not Supported"
	statusCodes[506] = "506 Variant Also Negotiates"
	statusCodes[507] = "507 Insufficient Storage"
	statusCodes[508] = "508 Loop Detected"
	statusCodes[509] = "509 Bandwidth Limit Exceeded"
	statusCodes[510] = "510 Not Extended"
	statusCodes[511] = "511 Network Authentication Required"
}