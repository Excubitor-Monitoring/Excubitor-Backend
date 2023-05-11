package db

import (
	"database/sql"
	"fmt"
)

type Reader struct {
	db *sql.DB
}

func GetReader() *Reader {
	return reader
}

func (reader *Reader) GetHistoryEntriesByTarget(target string) (History, error) {
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
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			logger.Error("Error upon closing rows:", err)
		}
	}(rows)

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
