package db

import (
	"database/sql"
	"github.com/Excubitor-Monitoring/Excubitor-Backend/internal/config"
	"github.com/Excubitor-Monitoring/Excubitor-Backend/internal/logging"
	_ "github.com/mattn/go-sqlite3"
	"sync"
)

var singletonOnce sync.Once

var logger logging.Logger

// singleton variables
var writer *Writer
var reader *Reader

const createStatement string = `
	CREATE TABLE IF NOT EXISTS history (
			time DATETIME NOT NULL,
			target TEXT NOT NULL,
			content TEXT NOT NULL,
			PRIMARY KEY (time, target)
	);
`

func InitDatabase() error {
	var err error

	singletonOnce.Do(func() {
		logger = logging.GetLogger()

		file := config.GetConfig().String("data.database_file")

		logger.Trace("Opening database connection!")
		db, err := sql.Open("sqlite3", file)
		if err != nil {
			return
		}

		logger.Trace("Creating table if it doesn't exist already.")
		if _, err := db.Exec(createStatement); err != nil {
			return
		}

		writer = &Writer{db}
		reader = &Reader{db}
	})

	if err != nil {
		return err
	}

	return nil
}
