package main

import (
	"net/http"
)

func fileServerHandler() http.Handler {
	fs := http.FileServer(http.Dir("."))
	return http.StripPrefix("/app", fs)
}
