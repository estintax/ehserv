package main

func findExt(ext string) bool {
	for i := 0; i < len(extEnum); i++ {
		if extEnum[i] == ext {
			return true
		}
	}

	return false
}
