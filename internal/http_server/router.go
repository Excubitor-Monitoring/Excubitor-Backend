package http_server

import (
	"net/http"
	"strings"
)

func Serve(w http.ResponseWriter, r *http.Request) {
	path := strings.Split(r.URL.Path, "/")
	length := len(path)

	var handler http.Handler

	switch {
	case length == 1 && path[0] == "info" && r.Method == "GET":
		handler = http.HandlerFunc(info)
	case length == 1 && path[0] == "ws":
		handler = queryAuth(http.HandlerFunc(wsInit))
	case length == 1 && path[0] == "auth":
		handler = http.HandlerFunc(handleAuthRequest)
	case length == 2 && path[0] == "auth" && path[1] == "refresh":
		handler = http.HandlerFunc(handleRefreshRequest)
	case path[0] == "static":
		handler = http.HandlerFunc(handleStaticFiles)
	default:
		http.NotFound(w, r)
		return
	}

	handler.ServeHTTP(w, r)
}
