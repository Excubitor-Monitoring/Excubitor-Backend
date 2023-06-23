package frontend

import (
	"embed"
	"errors"
	"fmt"
	"github.com/Excubitor-Monitoring/Excubitor-Backend/internal/http_server/helper"
	"github.com/Excubitor-Monitoring/Excubitor-Backend/internal/logging"
	"github.com/Excubitor-Monitoring/Excubitor-Backend/internal/plugins"
	"mime"
	"net/http"
	"path/filepath"
	"strconv"
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
		content, err := plugins.GetExternalFrontendFile(r.URL.Path)
		if err != nil {
			switch {
			case errors.Is(plugins.ErrFrontendComponentsNotProvided, err):
				helper.ReturnError(w, r, http.StatusNotImplemented, "Plugin does not implement frontend components!")
			case errors.Is(plugins.ErrFrontendComponentFileNotFound, err):
				helper.ReturnError(w, r, http.StatusNotFound, "Frontend component not found!")
			default:
				helper.ReturnError(w, r, http.StatusInternalServerError, "Unknown error!")
			}

			return
		}

		w.Header().Set("Content-Type", mime.TypeByExtension(filepath.Ext(path[length-1])))
		w.Header().Set("Content-Length", strconv.Itoa(len(content)))
		if _, err := w.Write(content); err != nil {
			logging.GetLogger().Error(fmt.Sprintf("Could not fulfill request for %s from %s!", r.URL.Path, r.RemoteAddr))
			return
		}
	default:
		http.NotFound(w, r)
	}

}
