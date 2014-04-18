package handler

import (
	"net/http"
)

func init() {
	http.Handle("/", http.FileServer(http.Dir("./static_files")))
}

