package http_server

import (
	"encoding/json"
	"fmt"
	ctx "github.com/Excubitor-Monitoring/Excubitor-Backend/internal/context"
	"github.com/Excubitor-Monitoring/Excubitor-Backend/internal/http_server/models"
	"github.com/Excubitor-Monitoring/Excubitor-Backend/internal/logging"
	"github.com/spf13/viper"
	"net/http"
)

var logger logging.Logger

func Start() error {
	host := viper.GetString("http.host")
	port := viper.GetInt("http.port")

	logger.Debug(fmt.Sprintf("Starting HTTP Server on port %d", port))

	mux := http.NewServeMux()
	mux.HandleFunc("/info", info)
	err := http.ListenAndServe(fmt.Sprintf("%s:%d", host, port), mux)
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

func init() {
	var err error
	logger, err = logging.GetLogger()
	if err != nil {
		panic(err)
	}
}
