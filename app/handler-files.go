package app

import (
	"net/http"
)

func FileServerHandler() http.Handler {
	fs := http.FileServer(http.Dir("./app"))
	return http.StripPrefix("/app", fs)
}
