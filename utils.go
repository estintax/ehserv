package main

import "time"

func findExt(ext string) bool {
	for i := 0; i < len(extEnum); i++ {
		if extEnum[i] == ext {
			return true
		}
	}

	return false
}

func getDate() string {
	curTime := time.Now().UTC()
	return curTime.Format("Mon, 02 Jan 2006 15:04:05 GMT")
}
