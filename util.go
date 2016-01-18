package main

import (
	"mime"
	"strings"
)

func getMimeType(filename string) string {
	idx := strings.LastIndex(filename, ".")
	if idx == -1 {
		return ""
	}
	return mime.TypeByExtension(filename[idx:])
}
