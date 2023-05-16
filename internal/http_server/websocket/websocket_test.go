package websocket

import (
	"compress/zlib"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/Excubitor-Monitoring/Excubitor-Backend/internal/config"
	ctx "github.com/Excubitor-Monitoring/Excubitor-Backend/internal/context"
	"github.com/Excubitor-Monitoring/Excubitor-Backend/internal/db"
	"github.com/gobwas/ws/wsutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestNewMessage(t *testing.T) {
	msg := NewMessage(UNSUB, "Some.Target.Address", "Some value")
	assert.Equal(t, UNSUB, msg.OpCode)
	assert.Equal(t, TargetAddress("Some.Target.Address"), msg.Target)
	assert.Equal(t, "Some value", msg.Value)
}

// Protocol integration tests

func TestSendMessage(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(1)

	server, client := net.Pipe()

	msg := NewMessage(UNSUB, "Some.Target.Address", "Some value")

	go func() {
		received, err := wsutil.ReadServerText(client)
		if err != nil {
			t.Error(err)
			return
		}

		expected, err := msg.Bytes()
		if err != nil {
			t.Error(err)
			return
		}

		assert.Equal(t, expected, received)
		wg.Done()
	}()

	if err := sendMessage(server, msg); err != nil {
		t.Error(err)
		return
	}

	wg.Wait()

}

func TestSUB(t *testing.T) {
	done := make(chan bool)
	timeout := time.After(1 * time.Second)

	server, client := net.Pipe()

	request := NewMessage(SUB, "Some.Target", "")

	go HandleWebsocket(server)

	bytes, err := request.Bytes()
	if err != nil {
		t.Error(err)
		return
	}

	if err := wsutil.WriteClientText(client, bytes); err != nil {
		t.Error(err)
		return
	}

	msg := NewMessage(REPLY, "Some.Target", "SomeValue")

	go func() {
		received, err := wsutil.ReadServerText(client)
		if err != nil {
			t.Error(err)
			return
		}

		expected, err := msg.Bytes()
		if err != nil {
			t.Error(err)
			return
		}

		assert.Equal(t, expected, received)
		done <- true
	}()

	go func() {
		for {
			ctx.GetContext().GetBroker().Publish("Some.Target", "SomeValue")
			time.Sleep(50 * time.Millisecond)
		}
	}()

	select {
	case <-timeout:
		t.Fatal("Test didn't finish in time...")
	case <-done:
	}
}

func TestGET(t *testing.T) {
	done := make(chan bool)
	timeout := time.After(1 * time.Second)

	server, client := net.Pipe()

	request := NewMessage(GET, "Some.Target", "")

	go HandleWebsocket(server)

	bytes, err := request.Bytes()
	if err != nil {
		t.Error(err)
		return
	}

	if err := wsutil.WriteClientText(client, bytes); err != nil {
		t.Error(err)
		return
	}

	msg := NewMessage(REPLY, "Some.Target", "SomeValue")

	go func() {
		received, err := wsutil.ReadServerText(client)
		if err != nil {
			t.Error(err)
			return
		}

		expected, err := msg.Bytes()
		if err != nil {
			t.Error(err)
			return
		}

		assert.Equal(t, expected, received)
		done <- true
	}()

	go func() {
		for {
			ctx.GetContext().GetBroker().Publish("Some.Target", "SomeValue")
			time.Sleep(50 * time.Millisecond)
		}
	}()

	select {
	case <-timeout:
		t.Fatal("Test didn't finish in time...")
	case <-done:
	}
}

func TestHIST(t *testing.T) {
	// SETUP TEST DATA

	var database *sql.DB
	database, err := sql.Open("sqlite3", config.GetConfig().String("data.database_file"))
	if err != nil {
		t.Error(err)
		return
	}

	transaction, err := database.Begin()
	if err != nil {
		t.Error(err)
		return
	}

	stmt, err := transaction.Prepare("INSERT INTO history (time, target, content) VALUES (?, ?, ?)")
	if err != nil {
		t.Error(err)
		return
	}

	for i := 0; i < 5; i++ {
		// COMPRESS TEST DATA VALUE FIELD
		buf := new(strings.Builder)

		messageValue := fmt.Sprintf("Message No. %d", i)

		writer := zlib.NewWriter(buf)
		if _, err := writer.Write([]byte(messageValue)); err != nil {
			t.Error(err)
			return
		}
		if err := writer.Close(); err != nil {
			t.Error(err)
			return
		}

		compressedValue := buf.String()

		_, err := stmt.Exec(time.Now().Add(-time.Duration(i)*time.Minute), "Some.Target", compressedValue)
		if err != nil {
			t.Error(err)
			return
		}

	}

	if err := transaction.Commit(); err != nil {
		t.Error(err)
		return
	}

	if err := database.Close(); err != nil {
		t.Error(err)
		return
	}

	// HELPER FUNCTION

	sendHISTRequest := func(client net.Conn, target string, params HistoryRequestParameters) (db.History, error) {
		paramJson, err := json.Marshal(params)
		if err != nil {
			return nil, err
		}

		request := NewMessage(HIST, "Some.Target", string(paramJson))
		requestBytes, err := request.Bytes()
		if err != nil {
			return nil, err
		}

		if err := wsutil.WriteClientText(client, requestBytes); err != nil {
			return nil, err
		}

		received, err := wsutil.ReadServerText(client)
		if err != nil {
			return nil, err
		}

		var reply Message
		if err := json.Unmarshal(received, &reply); err != nil {
			return nil, err
		}

		var history db.History
		if err := json.Unmarshal([]byte(reply.Value), &history); err != nil {
			return nil, err
		}

		return history, nil
	}

	// BEGIN ACTUAL TEST

	server, client := net.Pipe()
	go HandleWebsocket(server)

	history, err := sendHISTRequest(client, "Some.Target", HistoryRequestParameters{Until: time.Now()})
	if err != nil {
		t.Error(err)
	}

	require.Equal(t, 5, len(history))
	assert.Equal(t, "Message No. 4", history[0].Message.Value)
	assert.Equal(t, "Message No. 3", history[1].Message.Value)
	assert.Equal(t, "Message No. 2", history[2].Message.Value)
	assert.Equal(t, "Message No. 1", history[3].Message.Value)
	assert.Equal(t, "Message No. 0", history[4].Message.Value)
	assert.Equal(t, "Some.Target", history[0].Message.Target)
	assert.Equal(t, "Some.Target", history[1].Message.Target)
	assert.Equal(t, "Some.Target", history[2].Message.Target)
	assert.Equal(t, "Some.Target", history[3].Message.Target)
	assert.Equal(t, "Some.Target", history[4].Message.Target)

	history, err = sendHISTRequest(client, "Some.Target", HistoryRequestParameters{From: history[1].Timestamp, Until: history[3].Timestamp})
	if err != nil {
		t.Error(err)
		return
	}

	require.Equal(t, 3, len(history))
	assert.Equal(t, "Message No. 3", history[0].Message.Value)
	assert.Equal(t, "Message No. 2", history[1].Message.Value)
	assert.Equal(t, "Message No. 1", history[2].Message.Value)
	assert.Equal(t, "Some.Target", history[0].Message.Target)
	assert.Equal(t, "Some.Target", history[1].Message.Target)
	assert.Equal(t, "Some.Target", history[2].Message.Target)
}
