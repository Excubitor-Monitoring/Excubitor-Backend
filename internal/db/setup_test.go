package db

import (
	"github.com/Excubitor-Monitoring/Excubitor-Backend/internal/config"
	"github.com/Excubitor-Monitoring/Excubitor-Backend/internal/logging"
	"github.com/knadh/koanf/providers/confmap"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	err := config.GetConfig().Load(confmap.Provider(map[string]interface{}{
		"logging.log_level":  "TRACE",
		"logging.method":     "CONSOLE",
		"data.storage_time":  "720h",
		"data.purge_cycle":   "1h",
		"data.database_file": "history_test.db",
	}, "."), nil)
	if err != nil {
		panic(err)
	}

	if err := logging.InitLogging(); err != nil {
		panic(err)
	}

	code := m.Run()

	if err := os.Remove("history_test.db"); err != nil {
		panic(err)
	}

	os.Exit(code)
}

func clearDatabase() error {
	if _, err := GetWriter().db.Exec("DELETE FROM history WHERE true"); err != nil {
		return err
	}

	return nil
}
