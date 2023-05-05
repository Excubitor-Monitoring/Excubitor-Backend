package http_server

import (
	"encoding/json"
	"fmt"
	"github.com/Excubitor-Monitoring/Excubitor-Backend/internal/config"
	ctx "github.com/Excubitor-Monitoring/Excubitor-Backend/internal/context"
	"github.com/Excubitor-Monitoring/Excubitor-Backend/internal/http_server/models"
	"github.com/Excubitor-Monitoring/Excubitor-Backend/internal/logging"
	"github.com/gobwas/ws"
	"net/http"
)

var logger logging.Logger
var k = config.GetConfig()

func Start() error {
	host := k.String("http.host")
	port := k.Int("http.port")

	logger = logging.GetLogger()

	logger.Info(fmt.Sprintf("Starting HTTP Server on port %d", port))

	mux := http.NewServeMux()
	mux.HandleFunc("/info", info)
	mux.HandleFunc("/auth", handleAuthRequest)
	mux.HandleFunc("/auth/refresh", handleRefreshRequest)
	mux.Handle("/ws", queryAuth(http.HandlerFunc(wsInit)))

	cors := getCORSHandler()

	err := http.ListenAndServe(fmt.Sprintf("%s:%d", host, port), cors.Handler(mux))
	if err != nil {
		return err
	}

	return nil
}

func info(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		w.Header().Set("Content-Type", "application/json")
		jsonResult, err := json.Marshal(models.NewInfoResponse("PAM", ctx.GetContext().GetModules()))
		if err != nil {
			return
		}

		_, err = w.Write(jsonResult)
		if err != nil {
			return
		}
		break
	default:
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusMethodNotAllowed)

		jsonResult, err := json.Marshal(NewHTTPError(fmt.Sprintf("Method %s not allowed!", r.Method), r.RequestURI))
		if err != nil {
			return
		}

		_, err = w.Write(jsonResult)
		if err != nil {
			return
		}
	}
}

func wsInit(w http.ResponseWriter, r *http.Request) {
	conn, _, _, err := ws.UpgradeHTTP(r, w)
	if err != nil {
		logger.Error(fmt.Sprintf("Connection from %s couldn't be upgraded: %s", r.RemoteAddr, err))
	}

	handleWebsocket(conn)
}
