package db

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestWriter_AddHistoryEntry(t *testing.T) {
	if err := InitDatabase(); err != nil {
		t.Error(err)
		return
	}

	if err := clearDatabase(); err != nil {
		t.Error(err)
		return
	}

	writer := GetWriter()

	rows, err := writer.db.Query("SELECT count(*) FROM history")
	if err != nil {
		t.Error(err)
		return
	}

	rows.Next()

	var count int
	if err := rows.Scan(&count); err != nil {
		t.Error(err)
		return
	}

	if err := rows.Close(); err != nil {
		t.Error(err)
		return
	}

	assert.Equal(t, 0, count)

	if err := writer.AddHistoryEntry("SomeTarget", "SomeContent"); err != nil {
		t.Error(err)
		return
	}

	rows, err = writer.db.Query("SELECT count(*) FROM history")
	if err != nil {
		t.Error(err)
		return
	}

	rows.Next()

	if err := rows.Scan(&count); err != nil {
		t.Error(err)
		return
	}

	if err := rows.Close(); err != nil {
		t.Error(err)
		return
	}

	assert.Equal(t, 1, count)
}
