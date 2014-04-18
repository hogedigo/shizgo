package handler

import (
	"fmt"
	"net/http"
)

func init() {
	http.Handle("/handler", new(testHandler))
}

type testHandler struct {
}

func (h testHandler) ServeHTTP(w http.ResponseWriter, _ *http.Request) {
	fmt.Fprintf(w, "hello handler!")
}

