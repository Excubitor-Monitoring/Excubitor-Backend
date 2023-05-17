package db

import (
	"database/sql"
	"fmt"
	"time"
)

type Reader struct {
	db *sql.DB
}

func GetReader() *Reader {
	return reader
}

// GetHistoryEntries gets history all entries from the loaded database file
func (reader *Reader) GetHistoryEntries(target string) (History, error) {
	stmt, err := reader.db.Prepare(`
		SELECT * FROM history WHERE target = ?;
	`)
	if err != nil {
		return nil, err
	}

	return retrieveHistoryFromDB(stmt, target)
}

// GetHistoryEntriesThinned thins out the history data gathered
// from GetHistoryEntries to be at least as spread apart as defined in maxDensity.
func (reader *Reader) GetHistoryEntriesThinned(target string, maxDensity time.Duration) (History, error) {
	data, err := reader.GetHistoryEntries(target)
	if err != nil {
		return nil, err
	}

	return thinData(data, maxDensity), nil
}

// GetHistoryEntriesFrom gets History entries after "from"
func (reader *Reader) GetHistoryEntriesFrom(target string, from time.Time) (History, error) {
	return reader.GetHistoryEntriesFromUntil(target, from, time.Now())
}

// GetHistoryEntriesFromThinned thins out the history data gathered
// from GetHistoryEntriesFrom to be at least as spread apart as defined in maxDensity
func (reader *Reader) GetHistoryEntriesFromThinned(target string, from time.Time, maxDensity time.Duration) (History, error) {
	data, err := reader.GetHistoryEntriesFrom(target, from)
	if err != nil {
		return nil, err
	}

	return thinData(data, maxDensity), nil
}

// GetHistoryEntriesFromUntil gets History entries after "from" and before "until"
func (reader *Reader) GetHistoryEntriesFromUntil(target string, from time.Time, until time.Time) (History, error) {
	stmt, err := reader.db.Prepare(`
		SELECT * FROM history WHERE target = ? AND time >= ? AND time <= ?;
	`)
	if err != nil {
		return nil, err
	}

	return retrieveHistoryFromDB(stmt, target, from, until)
}

// GetHistoryEntriesFromUntilThinned thins out the history data gathered
// from GetHistoryEntriesFromUntil to be at least as speread apart as defined in maxDensity
func (reader *Reader) GetHistoryEntriesFromUntilThinned(target string, from time.Time, until time.Time, maxDensity time.Duration) (History, error) {
	data, err := reader.GetHistoryEntriesFromUntil(target, from, until)
	if err != nil {
		return nil, err
	}

	return thinData(data, maxDensity), nil
}

// retrieveHistoryFromDB calls the prepared statement stmt with the specified args and retrieves a History slice from the results.
func retrieveHistoryFromDB(stmt *sql.Stmt, args ...any) (History, error) {
	rows, err := stmt.Query(args...)
	if err != nil {
		return nil, err
	}

	data, err := collectHistoryMessages(rows)
	if err != nil {
		return nil, err
	}

	if err := rows.Close(); err != nil {
		logger.Error("Error upon closing rows:", err)
	}

	if err := stmt.Close(); err != nil {
		logger.Error("Error upon closing statement for reader:", err)
	}

	return data, nil
}

// thinData thins out the History specified to be at least as spread apart as specified in maxDensity
func thinData(history History, maxDensity time.Duration) History {
	if len(history) == 0 {
		return history
	}

	var newHistory History
	newHistory = append(newHistory, history[0])

	reference := history[0].Timestamp.Add(maxDensity)
	remaining := history[1:]

	for _, entry := range remaining {
		if entry.Timestamp.After(reference) || entry.Timestamp.Equal(reference) {
			newHistory = append(newHistory, entry)
			reference = entry.Timestamp.Add(maxDensity)
		}
	}

	return newHistory
}

// collectHistoryMessages constructs a History from database rows.
func collectHistoryMessages(rows *sql.Rows) (History, error) {
	data := History{}
	for rows.Next() {
		message := HistoryMessage{}
		err := rows.Scan(&message.Timestamp, &message.Message.Target, &message.Message.Value)
		if err != nil {
			return nil, err
		}

		decompressedValue, err := decompress(message.Message.Value)
		if err != nil {
			logger.Error(fmt.Sprintf("Could not decompress value of timestamp %s of target %s! Reason: %s", message.Timestamp.UTC().String(), message.Message.Target, err))
			continue
		}

		message.Message.Value = decompressedValue

		data = append(data, message)
	}

	return data, nil
}
