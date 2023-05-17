package http_server

import (
	"github.com/Excubitor-Monitoring/Excubitor-Backend/internal/config"
	"github.com/Excubitor-Monitoring/Excubitor-Backend/internal/logging"
	"github.com/knadh/koanf/providers/confmap"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	err := config.GetConfig().Load(confmap.Provider(map[string]interface{}{
		"logging.log_level":         "TRACE",
		"logging.method":            "CONSOLE",
		"http.host":                 "0.0.0.0",
		"http.port":                 "8080",
		"http.cors.allowed_origins": []string{"*"},
		"http.cors.allowed_methods": []string{"GET", "POST"},
		"http.cors.allowed_headers": []string{"Origin", "Content-Type", "Authorization"},
		"data.module_clock":         "5s",
		"data.storage_time":         "30d",
	}, "."), nil)
	if err != nil {
		panic(err)
	}

	if err := logging.InitLogging(); err != nil {
		panic(err)
	}

	code := m.Run()
	os.Exit(code)
}
