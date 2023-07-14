package db

import (
	"database/sql"
	"fmt"
	"github.com/Excubitor-Monitoring/Excubitor-Backend/internal/config"
	"github.com/Excubitor-Monitoring/Excubitor-Backend/internal/logging"
	_ "github.com/mattn/go-sqlite3"
	"sync"
	"time"
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

// InitDatabase initializes the database connection and starts all recurring jobs on the database.
func InitDatabase() error {
	var err error

	singletonOnce.Do(func() {
		logger = logging.GetLogger()

		file := config.GetConfig().String("data.database_file")

		logger.Trace("Opening database connection!")

		var db *sql.DB
		db, err = sql.Open("sqlite3", file)
		if err != nil {
			return
		}

		logger.Trace("Creating table if it doesn't exist already.")
		_, err = db.Exec(createStatement)
		if err != nil {
			return
		}

		writer = &Writer{db}
		reader = &Reader{db}

		err = vacuumDB(db)
		if err != nil {
			return
		}

		err = startPurgeCycle(db)
		if err != nil {
			return
		}
	})

	if err != nil {
		return err
	}

	return nil
}

// startPurgeCycle starts the recurring job of purging all old database entries.
func startPurgeCycle(db *sql.DB) error {
	purgeCycleString := config.GetConfig().String("data.purge_cycle")
	purgeCycle, err := time.ParseDuration(purgeCycleString)
	if err != nil {
		return err
	}

	logger.Trace("Starting purge cycle...")

	go func() {
		for {
			if err := purgeOldEntries(db); err != nil {
				logger.Error("Could not purge old database entries! Reason:", err.Error())
			}

			time.Sleep(purgeCycle)
		}
	}()

	return nil
}

// purgeOldEntries purges all old database entries.
func purgeOldEntries(db *sql.DB) error {
	storageTimeString := config.GetConfig().String("data.storage_time")
	storageTime, err := time.ParseDuration(storageTimeString)
	if err != nil {
		return err
	}

	stmt, err := db.Prepare(`DELETE FROM history WHERE time < ?`)
	if err != nil {
		return err
	}

	logger.Debug("Purging database...")

	result, err := stmt.Exec(time.Now().Add(-storageTime))
	if err != nil {
		_ = stmt.Close()
		return err
	}

	if err := stmt.Close(); err != nil {
		logger.Error("Error on closing statement for purging database entries:", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		logger.Error("Couldn't determine how many rows were deleted on purge!")
		return nil
	}

	logger.Debug(fmt.Sprintf("Deleted %d rows on purge.", rowsAffected))

	return nil
}

// vacuumDB executes a VACUUM statement on the sqlite database.
// This shrinks the database's size so that it does not become unnecessarily large.
func vacuumDB(db *sql.DB) error {
	logger.Debug("Vacuuming database...")

	stmt, err := db.Prepare(`VACUUM`)
	if err != nil {
		return err
	}

	if _, err := stmt.Exec(); err != nil {
		return err
	}

	if err := stmt.Close(); err != nil {
		logger.Error("Error on closing statement for vacuuming database:", err)
	}

	return nil
}
