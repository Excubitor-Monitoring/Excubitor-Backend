package http_server

import (
	"github.com/Excubitor-Monitoring/Excubitor-Backend/internal/frontend"
	"net/http"
	"strings"
)

func Serve(w http.ResponseWriter, r *http.Request) {
	path := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
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
	case path[0] == "static" && r.Method == "GET":
		handler = http.HandlerFunc(frontend.StaticFileServer)
	default:
		http.NotFound(w, r)
		return
	}

	handler.ServeHTTP(w, r)
}
