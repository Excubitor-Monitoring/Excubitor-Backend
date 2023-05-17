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

func TestHandleWebsocketBadRequest(t *testing.T) {
	server, client := net.Pipe()

	go HandleWebsocket(server)

	if err := wsutil.WriteClientText(client, []byte("Invalid message!")); err != nil {
		t.Error(err)
		return
	}

	replyBytes, err := wsutil.ReadServerText(client)
	if err != nil {
		t.Error(err)
		return
	}

	assertMessage(t, ERR, string(GetEmptyTarget()), "Bad Request!", replyBytes)
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

func TestREPLY(t *testing.T) {
	server, client := net.Pipe()

	go HandleWebsocket(server)

	request, err := Message{REPLY, "Some.Target", "Some value!"}.Bytes()
	if err != nil {
		t.Error(err)
		return
	}

	if err := wsutil.WriteClientText(client, request); err != nil {
		t.Error(err)
		return
	}

	replyBytes, err := wsutil.ReadServerText(client)
	if err != nil {
		t.Error(err)
		return
	}

	assertMessage(t, ERR, "Some.Target", "Clients may not send messages of the type REPLY!", replyBytes)
}

func TestERR(t *testing.T) {
	server, client := net.Pipe()

	go HandleWebsocket(server)

	request, err := Message{ERR, "Some.Target", "Some value!"}.Bytes()
	if err != nil {
		t.Error(err)
		return
	}

	if err := wsutil.WriteClientText(client, request); err != nil {
		t.Error(err)
		return
	}

	replyBytes, err := wsutil.ReadServerText(client)
	if err != nil {
		t.Error(err)
		return
	}

	assertMessage(t, ERR, "Some.Target", "Clients may not send messages of the type ERR!", replyBytes)
}

func TestUnsupportedOption(t *testing.T) {
	server, client := net.Pipe()

	go HandleWebsocket(server)

	request, err := Message{"UNSUPPORTED", "Some.Target", "Some value!"}.Bytes()
	if err != nil {
		t.Error(err)
		return
	}

	if err := wsutil.WriteClientText(client, request); err != nil {
		t.Error(err)
		return
	}

	replyBytes, err := wsutil.ReadServerText(client)
	if err != nil {
		t.Error(err)
		return
	}

	assertMessage(t, ERR, "Some.Target", "Unsupported Operation UNSUPPORTED!", replyBytes)
}

func TestSUB(t *testing.T) {
	done := make(chan bool)
	timeout := time.After(1 * time.Second)

	server, client := net.Pipe()

	request := NewMessage(SUB, "Some.Target.SUB", "")

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

	msg := NewMessage(REPLY, "Some.Target.SUB", "Some Value!")

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

	quit := make(chan bool)

	go func() {
		for {
			select {
			case <-quit:
				break
			default:
				ctx.GetContext().GetBroker().Publish("Some.Target.SUB", "Some Value!")
				time.Sleep(50 * time.Millisecond)
			}
		}
	}()

	select {
	case <-timeout:
		quit <- true
		t.Fatal("Test didn't finish in time...")
	case <-done:
		quit <- true
	}
}

func TestUNSUB(t *testing.T) {
	server, client := net.Pipe()

	go HandleWebsocket(server)

	request, err := Message{SUB, "Some.Target.UNSUB", ""}.Bytes()
	if err != nil {
		t.Error(err)
		return
	}

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		response, err := wsutil.ReadServerText(client)
		if err != nil {
			t.Error(err)
			return
		}

		assertMessage(t, REPLY, "Some.Target.UNSUB", "Some Value!", response)
		wg.Done()
	}()

	if err := wsutil.WriteClientText(client, request); err != nil {
		t.Error(err)
		return
	}

	quit := make(chan bool)

	go func() {
		for {
			select {
			case <-quit:
				break
			default:
				ctx.GetContext().GetBroker().Publish("Some.Target.UNSUB", "Some Value!")
				time.Sleep(50 * time.Millisecond)
			}
		}
	}()

	wg.Wait()

	unsubRequest, err := Message{UNSUB, "Some.Target.UNSUB", ""}.Bytes()
	if err != nil {
		t.Error(err)
		return
	}

	if err := wsutil.WriteClientText(client, unsubRequest); err != nil {
		t.Error(err)
		return
	}

	fail := make(chan bool)
	timeout := time.After(100 * time.Millisecond)

	go func() {
		_, err := wsutil.ReadServerText(client)
		if err != nil {
			t.Error(err)
			return
		}

		fail <- true
	}()

	select {
	case <-fail:
		t.Fatal("Test failed! Received message after UNSUB request!")
	case <-timeout:
		t.Log("Timeout on listening for response. Test finished successfully!")
	}

	quit <- true
}

func TestGET(t *testing.T) {
	done := make(chan bool)
	timeout := time.After(1 * time.Second)

	server, client := net.Pipe()

	request := NewMessage(GET, "Some.Target.GET", "")

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

	msg := NewMessage(REPLY, "Some.Target.GET", "Some Value!")

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

	quit := make(chan bool)

	go func() {
		go func() {
			for {
				select {
				case <-quit:
					break
				default:
					ctx.GetContext().GetBroker().Publish("Some.Target.GET", "Some Value!")
					time.Sleep(50 * time.Millisecond)
				}
			}
		}()
	}()

	select {
	case <-timeout:
		t.Fatal("Test didn't finish in time...")
	case <-done:
	}

	quit <- true
}

func TestHIST(t *testing.T) {
	// SETUP TEST DATA

	var database *sql.DB
	database, err := sql.Open("sqlite3", config.GetConfig().String("data.database_file"))
	if err != nil {
		t.Error(err)
		return
	}

	_, err = database.Exec(`DELETE FROM history WHERE true`)
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

	reference := time.Now()

	for i := 0; i < 5; i++ {
		// COMPRESS TEST DATA VALUE FIELD
		buf := new(strings.Builder)

		messageValue := fmt.Sprintf("Message No. %d", 4-i)

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

		_, err := stmt.Exec(reference.Add(-time.Duration(i)*time.Minute), "Some.Target", compressedValue)
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

		reply, err := decodeMessage(received)

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
	assert.Equal(t, "Message No. 0", history[0].Message.Value)
	assert.Equal(t, "Message No. 1", history[1].Message.Value)
	assert.Equal(t, "Message No. 2", history[2].Message.Value)
	assert.Equal(t, "Message No. 3", history[3].Message.Value)
	assert.Equal(t, "Message No. 4", history[4].Message.Value)
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
	assert.Equal(t, "Message No. 1", history[0].Message.Value)
	assert.Equal(t, "Message No. 2", history[1].Message.Value)
	assert.Equal(t, "Message No. 3", history[2].Message.Value)
	assert.Equal(t, "Some.Target", history[0].Message.Target)
	assert.Equal(t, "Some.Target", history[1].Message.Target)
	assert.Equal(t, "Some.Target", history[2].Message.Target)

	history, err = sendHISTRequest(client, "Some.Target", HistoryRequestParameters{Until: time.Now(), MaxDensity: 2 * time.Minute})

	require.Equal(t, 3, len(history))
	assert.Equal(t, "Message No. 0", history[0].Message.Value)
	assert.Equal(t, "Message No. 2", history[1].Message.Value)
	assert.Equal(t, "Message No. 4", history[2].Message.Value)
}

func TestHISTNegative(t *testing.T) {
	server, client := net.Pipe()
	go HandleWebsocket(server)

	request := Message{
		OpCode: HIST,
		Target: "Some.Target",
		Value:  "abcdefghijklmnopqrstuvwxyz",
	}

	requestBytes, err := json.Marshal(request)
	if err != nil {
		t.Error(err)
		return
	}

	if err := wsutil.WriteClientText(client, requestBytes); err != nil {
		t.Error(err)
		return
	}

	received, err := wsutil.ReadServerText(client)
	if err != nil {
		t.Error(err)
		return
	}

	assertMessage(t, ERR, "Some.Target", "Bad parameters!", received)

}

func decodeMessage(messageJSON []byte) (Message, error) {
	var output Message
	if err := json.Unmarshal(messageJSON, &output); err != nil {
		return Message{}, err
	}

	return output, nil
}

func assertMessage(t *testing.T, expectedOpCode OpCode, expectedTarget string, expectedValue string, receivedWebsocketMessage []byte) {
	msg, err := decodeMessage(receivedWebsocketMessage)
	if err != nil {
		t.Error(err)
		return
	}

	assert.Equal(t, expectedOpCode, msg.OpCode)
	assert.Equal(t, TargetAddress(expectedTarget), msg.Target)
	assert.Equal(t, expectedValue, msg.Value)
}
