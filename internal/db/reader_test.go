package db

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestReader_GetHistoryEntriesByTarget(t *testing.T) {
	if err := InitDatabase(); err != nil {
		t.Error(err)
		return
	}

	if err := clearDatabase(); err != nil {
		t.Error(err)
		return
	}

	reader := GetReader()

	compressedContent1, err := compress("Some content1")
	if err != nil {
		t.Error(err)
	}

	compressedContent2, err := compress("Some content2")
	if err != nil {
		t.Error(err)
	}

	compressedContent3, err := compress("Some content3")
	if err != nil {
		t.Error(err)
	}

	stmt, err := reader.db.Prepare("INSERT INTO history (time, target, content) VALUES (?, ?, ?)")
	if err != nil {
		t.Error(err)
		return
	}

	if _, err := stmt.Exec(time.Now(), "Target1", compressedContent1); err != nil {
		t.Error(err)
		return
	}

	if _, err := stmt.Exec(time.Now(), "Target1", compressedContent2); err != nil {
		t.Error(err)
		return
	}

	if _, err := stmt.Exec(time.Now(), "Target2", compressedContent3); err != nil {
		t.Error(err)
		return
	}

	if _, err := stmt.Exec(time.Now(), "Target3", compressedContent3); err != nil {
		t.Error(err)
		return
	}

	history, err := reader.GetHistoryEntries("Target1")
	if err != nil {
		return
	}

	assert.Equal(t, 2, len(history))

	oneContent1 := false
	oneContent2 := false

	for _, entry := range history {
		assert.Equal(t, "Target1", entry.Message.Target)
		assert.NotEqual(t, "Some content3", entry.Message.Value)

		if entry.Message.Value == "Some content1" {
			oneContent1 = true
		}

		if entry.Message.Value == "Some content2" {
			oneContent2 = true
		}
	}

	assert.True(t, oneContent1)
	assert.True(t, oneContent2)
}

func TestThinData(t *testing.T) {
	type testParameters struct {
		originalDuration time.Duration
		thinnedDuration  time.Duration
		originalLength   int
		thinnedLength    int
	}

	for _, params := range []testParameters{
		{
			originalDuration: 5 * time.Minute,
			thinnedDuration:  6 * time.Minute,
			originalLength:   7,
			thinnedLength:    4,
		},
		{
			originalDuration: 2 * time.Minute,
			thinnedDuration:  6 * time.Minute,
			originalLength:   8,
			thinnedLength:    3,
		},
	} {
		t.Run(
			fmt.Sprintf(
				"TestParameters{ originalDuration: %s, thinnedDuration: %s, originalLength: %d, thinnedLength: %d}",
				params.originalDuration, params.thinnedDuration, params.originalLength, params.thinnedLength,
			),
			func(t *testing.T) {
				history := History{}

				reference := time.Now()

				for i := 0; i < params.originalLength; i++ {
					history = append(history, HistoryMessage{
						Timestamp: reference.Add(time.Duration(i) * params.originalDuration),
					})
				}

				newHistory := thinData(history, params.thinnedDuration)

				assert.Equal(t, params.thinnedLength, len(newHistory))
			})
	}
}
