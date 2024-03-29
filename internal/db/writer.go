package db

import (
	"database/sql"
	"time"
)

type Writer struct {
	db *sql.DB
}

func GetWriter() *Writer {
	return writer
}

// AddHistoryEntry adds an entry to the history table.
func (writer *Writer) AddHistoryEntry(target string, content string) error {
	stmt, err := writer.db.Prepare(`
		INSERT INTO history (time, target, content) VALUES (?, ?, ?);
	`)
	if err != nil {
		return err
	}

	compressedValue, err := compress(content)
	if err != nil {
		return err
	}

	_, err = stmt.Exec(time.Now().UTC(), target, compressedValue)
	if err != nil {
		return err
	}

	if err := stmt.Close(); err != nil {
		logger.Error("Error on closing statement for writer:", err)
	}

	return nil
}
