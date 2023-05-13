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

func (reader *Reader) GetHistoryEntries(target string) (History, error) {
	stmt, err := reader.db.Prepare(`
		SELECT * FROM history WHERE target = ?;
	`)
	if err != nil {
		return nil, err
	}

	rows, err := stmt.Query(target)
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

func (reader *Reader) GetHistoryEntriesFromUntil(target string, from time.Time, until time.Time) (History, error) {
	stmt, err := reader.db.Prepare(`
		SELECT * FROM history WHERE target = ? AND time > ? AND time < ?;
	`)
	if err != nil {
		return nil, err
	}

	rows, err := stmt.Query(target, from, until)
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

func (reader *Reader) GetHistoryEntriesFrom(target string, from time.Time) (History, error) {
	return reader.GetHistoryEntriesFromUntil(target, from, time.Now())
}

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
			logger.Error(fmt.Sprintf("Could not decompress value of timestamp %s of target %s!", message.Timestamp.UTC().String(), message.Message.Target))
			continue
		}

		message.Message.Value = decompressedValue

		data = append(data, message)
	}

	return data, nil
}
