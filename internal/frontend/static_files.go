package frontend

import (
	"embed"
	"github.com/Excubitor-Monitoring/Excubitor-Backend/internal/http_server"
	"net/http"
	"strings"
)

//go:embed static
var embeddedFiles embed.FS

func StaticFileServer(w http.ResponseWriter, r *http.Request) {
	if strings.HasSuffix(r.URL.Path, "/") {
		http.NotFound(w, r)
		return
	}

	path := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	length := len(path)

	switch {
	case length > 1 && path[1] == "internal":
		http.FileServer(http.FS(embeddedFiles)).ServeHTTP(w, r)
	case length > 1 && path[1] == "external":
		http_server.ReturnError(w, r, http.StatusNotImplemented, "External modules are not implemented yet!")
	default:
		http.NotFound(w, r)
	}

}
