package http_server

import (
	"encoding/json"
	"fmt"
	ctx "github.com/Excubitor-Monitoring/Excubitor-Backend/internal/context"
	"github.com/Excubitor-Monitoring/Excubitor-Backend/internal/http_server/models"
	"github.com/Excubitor-Monitoring/Excubitor-Backend/internal/logging"
	"github.com/gobwas/ws"
	"github.com/spf13/viper"
	"net/http"
)

var logger logging.Logger

func Start() error {
	host := viper.GetString("http.host")
	port := viper.GetInt("http.port")

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
