package websocket

import (
	"errors"
	"github.com/Excubitor-Monitoring/Excubitor-Backend/internal/config"
	ctx "github.com/Excubitor-Monitoring/Excubitor-Backend/internal/context"
	"github.com/Excubitor-Monitoring/Excubitor-Backend/internal/db"
	"github.com/Excubitor-Monitoring/Excubitor-Backend/internal/logging"
	"github.com/Excubitor-Monitoring/Excubitor-Backend/internal/pubsub"
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

	if err := db.InitDatabase(); err != nil {
		panic(err)
	}

	ctx.GetContext().RegisterBroker(pubsub.NewBroker())

	code := m.Run()

	if _, err := os.Stat("history_test.db"); !errors.Is(err, os.ErrNotExist) {
		if err := os.Remove("history_test.db"); err != nil {
			panic(err)
		}
	}

	os.Exit(code)
}
