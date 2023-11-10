package main

import (
	"fmt"
	"io"
	"net/http"
)

func Handler(rw http.ResponseWriter, req *http.Request) {
	rw.Header().Set("X-Custom-Header", "Test")
	rw.WriteHeader(200)
	name := req.URL.Query().Get("name")
	_, _ = io.WriteString(rw, fmt.Sprintf("Hello, %s!", name))
}
