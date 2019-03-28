package main

func getContentType(format string) string {
	var contentType string
	switch format {
	case "html", "htm": contentType = "text/html"
	case "ogg": contentType = "application/ogg"
	case "mp3": contentType = "audio/mpeg"
	case "jpg": contentType = "image/jpeg"
	case "png": contentType = "image/png"
	case "gif": contentType = "image/gif"
	case "svg": contentType = "image/svg+xml"
	case "xml": contentType = "application/xml"
	case "js": contentType = "application/javascript"
	case "json": contentType = "application/json"
	case "css": contentType = "text/css"
	default: contentType = defaultMime
	}

	return contentType
}