package frontend

import (
	"embed"
	"net/http"
	"strings"
)

//go:embed static
var embeddedFiles embed.FS

func StaticFileServer(w http.ResponseWriter, r *http.Request) {
	path := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	length := len(path)

	switch {
	case length > 1 && path[1] == "internal":
		http.FileServer(http.FS(embeddedFiles)).ServeHTTP(w, r)
	case length > 1 && path[1] == "external":
		w.WriteHeader(http.StatusNotImplemented)
		return
	default:
		http.NotFound(w, r)
	}

}
