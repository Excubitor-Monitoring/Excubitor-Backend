package http_server

import (
	"fmt"
	"github.com/Excubitor-Monitoring/Excubitor-Backend/internal/frontend"
	"github.com/Excubitor-Monitoring/Excubitor-Backend/internal/http_server/helper"
	"net/http"
	"strings"
)

// Serve routes the incoming requests by URL path and path length
func Serve(w http.ResponseWriter, r *http.Request) {
	var remoteAddress string

	if remoteAddress = r.Header.Get("X-Forwarded-For"); remoteAddress == "" {
		remoteAddress = r.RemoteAddr
	}

	logger.Trace(fmt.Sprintf("[%s]: %s - %s", remoteAddress, r.Method, r.URL.Path))

	path := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	length := len(path)

	var handler http.Handler

	switch {
	case length == 1 && path[0] == "info" && r.Method == "GET":
		logger.Trace(fmt.Sprintf("[%s]: %s - %s -> info endpoint", remoteAddress, r.URL.Path, remoteAddress))
		handler = http.HandlerFunc(info)
	case length == 1 && path[0] == "ws":
		logger.Trace(fmt.Sprintf("[%s]: %s - %s -> ws endpoint, auth pending", remoteAddress, r.URL.Path, remoteAddress))
		handler = queryAuth(http.HandlerFunc(wsInit))
	case length == 1 && path[0] == "auth":
		logger.Trace(fmt.Sprintf("[%s]: %s - %s -> auth request endpoint", remoteAddress, r.URL.Path, remoteAddress))
		handler = http.HandlerFunc(handleAuthRequest)
	case length == 2 && path[0] == "auth" && path[1] == "refresh":
		logger.Trace(fmt.Sprintf("[%s]: %s - %s -> auth refresh endpoint", remoteAddress, r.URL.Path, remoteAddress))
		handler = http.HandlerFunc(handleRefreshRequest)
	case path[0] == "static" && r.Method == "GET":
		logger.Trace(fmt.Sprintf("[%s]: %s - %s -> static fileserver", remoteAddress, r.URL.Path, remoteAddress))
		handler = http.HandlerFunc(frontend.StaticFileServer)
	default:
		logger.Trace(fmt.Sprintf("[%s]: %s - %s -> NOT FOUND", remoteAddress, r.URL.Path, remoteAddress))
		helper.ReturnError(w, r, 404, fmt.Sprintf("Couldn't find requested resource %s!", r.URL.Path))
		http.NotFound(w, r)
		return
	}

	handler.ServeHTTP(w, r)
}
