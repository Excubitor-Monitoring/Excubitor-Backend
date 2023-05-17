package db

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestInitDatabaseCreateTable(t *testing.T) {
	err := InitDatabase()
	if err != nil {
		t.Error(err)
		return
	}

	query, err := GetReader().db.Query("SELECT count(*) FROM sqlite_master WHERE type='table' AND name='history';")
	if err != nil {
		t.Error(err)
		return
	}

	if query.Next() {
		var count int
		err := query.Scan(&count)
		if err != nil {
			t.Error(err)
			return
		}

		assert.Equal(t, 1, count)
	}

	if err := query.Close(); err != nil {
		t.Error(err)
		return
	}
}

func TestPurgeDatabaseEntries(t *testing.T) {
	err := InitDatabase()
	if err != nil {
		t.Error(err)
		return
	}

	if err := clearDatabase(); err != nil {
		t.Error(err)
		return
	}

	stmt, err := GetWriter().db.Prepare("INSERT INTO history (time, target, content) VALUES (?, ?, ?)")
	if err != nil {
		t.Error(err)
		return
	}

	_, err = stmt.Exec(time.Now().Add(-720*time.Hour), "Some.Target", "Some content")
	if err != nil {
		t.Error(err)
		return
	}

	_, err = stmt.Exec(time.Now().Add(-9600*time.Hour), "Some.Target", "Some content")
	if err != nil {
		t.Error(err)
		return
	}

	timeLower := time.Now().Add(-710 * time.Hour)

	_, err = stmt.Exec(timeLower, "Some.Target", "Some content")
	if err != nil {
		t.Error(err)
		return
	}

	if err := stmt.Close(); err != nil {
		t.Error(err)
		return
	}

	query, err := GetWriter().db.Query("SELECT count(*) FROM history")
	query.Next()

	var count int
	err = query.Scan(&count)
	if err != nil {
		t.Error(err)
		return
	}

	if err := query.Close(); err != nil {
		t.Error(err)
		return
	}

	assert.Equal(t, 3, count)

	if err := purgeOldEntries(GetWriter().db); err != nil {
		t.Error(err)
		return
	}

	query, err = GetWriter().db.Query("SELECT count(*) FROM history")
	query.Next()

	err = query.Scan(&count)
	if err != nil {
		t.Error(err)
		return
	}

	if err := query.Close(); err != nil {
		t.Error(err)
		return
	}

	assert.Equal(t, 1, count)

	query, err = GetWriter().db.Query("SELECT time, target, content FROM history")
	query.Next()

	var timestamp time.Time
	var target string
	var content string

	if err := query.Scan(&timestamp, &target, &content); err != nil {
		t.Error(err)
		return
	}

	if err := query.Close(); err != nil {
		t.Error(err)
		return
	}

	assert.True(t, timeLower.Equal(timestamp))
	assert.Equal(t, "Some.Target", target)
	assert.Equal(t, "Some content", content)

}
