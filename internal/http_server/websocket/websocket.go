package websocket

import (
	"encoding/json"
	"errors"
	"fmt"
	ctx "github.com/Excubitor-Monitoring/Excubitor-Backend/internal/context"
	"github.com/Excubitor-Monitoring/Excubitor-Backend/internal/db"
	"github.com/Excubitor-Monitoring/Excubitor-Backend/internal/logging"
	"github.com/Excubitor-Monitoring/Excubitor-Backend/internal/pubsub"
	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"net"
	"sync"
	"time"
)

var logger logging.Logger

var FatalWebsocketError error = errors.New("fatal websocket error")

func HandleWebsocket(conn net.Conn) {
	var err error
	clientAddress := conn.RemoteAddr()

	logger = logging.GetLogger()

	defer func(conn net.Conn) {
		logger.Debug(fmt.Sprintf("Closing connection from %s", clientAddress))

		err := conn.Close()
		if err != nil {
			logger.Error(fmt.Sprintf("Couldn't close connection from %s", clientAddress))
		}
	}(conn)

	// Set up subscriber

	broker := ctx.GetContext().GetBroker()
	subscriber := broker.AddSubscriber()
	defer func(subscriber *pubsub.Subscriber) {
		logger.Trace(fmt.Sprintf("Destructing subscriber associated with connection from %s", clientAddress))
		subscriber.Destruct()
	}(subscriber)

	go subscriber.Listen(func(m *pubsub.Message) {
		logger.Trace(fmt.Sprintf("Sending message from %s to connection from %s", m.GetMonitor(), clientAddress))

		if err = sendMessage(conn, NewMessage(REPLY, TargetAddress(m.GetMonitor()), m.GetMessageBody())); err != nil {
			return
		}
	})

	for {
		// Receiving message

		msg, op, err := wsutil.ReadClientData(conn)

		if err != nil {
			logger.Warn(fmt.Sprintf("Can't receive message from %s! Aborting connection...", clientAddress))
			return
		}

		if op == ws.OpClose {
			logger.Debug(fmt.Sprintf("Client from %s closes connection.", clientAddress))

			if err := conn.Close(); err != nil {
				logger.Error(fmt.Sprintf("Couldn't close connection from %s after client closed websocket!", clientAddress))
			}

			return
		}

		logger.Trace(fmt.Sprintf("Received message from %s: %s", clientAddress, string(msg)))

		// Decoding message

		content := &Message{}

		if err = json.Unmarshal(msg, content); err != nil {
			logger.Warn(fmt.Sprintf("Can't decode message from %s with reason %s! Dropping request...", clientAddress, err))

			if err := sendMessage(conn, NewMessage(ERR, GetEmptyTarget(), "Bad Request!")); err != nil {
				return
			}

			continue
		}

		switch content.OpCode {
		case GET:
			if err := handleGET(conn, content); err != nil {
				if errors.Is(FatalWebsocketError, err) {
					logger.Error(fmt.Sprintf("A fatal websocket error occurred. Forcefully aborting connection to %s!", clientAddress))
					return
				}
			}
		case SUB:
			broker.Subscribe(subscriber, string(content.Target))
			logger.Trace(fmt.Sprintf("Client %s subscribed to monitor %s.", clientAddress, content.Target))
		case UNSUB:
			broker.Unsubscribe(subscriber, string(content.Target))
			logger.Trace(fmt.Sprintf("Client %s unsubscribed from monitor %s.", clientAddress, content.Target))
		case HIST:
			if err := handleHIST(conn, content); err != nil {
				if errors.Is(FatalWebsocketError, err) {
					logger.Error(fmt.Sprintf("A fatal websocket error occurred. Forcefully aborting connection to %s!", clientAddress))
					return
				}
			}
		case REPLY:
			if err = sendMessage(conn, NewMessage(ERR, content.Target, "Clients may not send messages of the type REPLY!")); err != nil {
				return
			}
		case ERR:
			if err = sendMessage(conn, NewMessage(ERR, content.Target, "Clients may not send messages of the type ERR!")); err != nil {
				return
			}
		default:
			if err = sendMessage(conn, NewMessage(ERR, content.Target, fmt.Sprintf("Unsupported Operation %s!", content.OpCode))); err != nil {
				return
			}
		}
	}
}

func handleGET(conn net.Conn, content *Message) error {
	broker := ctx.GetContext().GetBroker()
	temporarySubscriber := broker.AddSubscriber()

	var receiveOnce sync.Once

	logger.Trace(fmt.Sprintf("Added temporary subscriber to fulfill GET request from %s on monitor %s.", conn.RemoteAddr(), content.Target))

	go temporarySubscriber.Listen(func(m *pubsub.Message) {
		receiveOnce.Do(func() {
			logger.Trace(fmt.Sprintf("Sending single message from %s to connection from %s", m.GetMonitor(), conn.RemoteAddr()))

			broker.Unsubscribe(temporarySubscriber, string(content.Target))
			defer temporarySubscriber.Destruct()

			if err := sendMessage(conn, NewMessage(REPLY, TargetAddress(m.GetMonitor()), m.GetMessageBody())); err != nil {
				return
			}
		})
	})

	broker.Subscribe(temporarySubscriber, string(content.Target))

	return nil
}

func handleHIST(conn net.Conn, content *Message) error {
	clientAddress := conn.RemoteAddr()

	logger.Debug("Client", clientAddress, "requested history of monitor", content.Target)

	params := &HistoryRequestParameters{
		From:       time.Time{},
		Until:      time.Now(),
		MaxDensity: "",
	}

	if content.Value != "" {
		if err := json.Unmarshal([]byte(content.Value), params); err != nil {
			logger.Error(fmt.Sprintf("Could not decode the history request parameters from %s. Reason: %s", clientAddress, err))

			if err := sendMessage(conn, NewMessage(ERR, content.Target, "Bad parameters!")); err != nil {
				return fmt.Errorf("%w: %s", FatalWebsocketError, err)
			}

			return err
		}
	}

	logger.Trace(fmt.Sprintf("Parsed HistoryRequestParameters { From: %s, Until: %s, MaxDensity: %s } from %s", params.From, params.Until, params.MaxDensity, clientAddress))

	var history db.History
	reader := db.GetReader()

	if params.MaxDensity != "" {
		var maxDensity time.Duration
		switch duration := params.MaxDensity.(type) {
		case float64:
			maxDensity = time.Duration(duration)
		case string:
			var err error
			maxDensity, err = time.ParseDuration(duration)
			if err != nil {
				return err
			}
		}

		var err error
		history, err = reader.GetHistoryEntriesFromUntilThinned(string(content.Target), params.From, params.Until, maxDensity)
		if err != nil {
			logger.Error(fmt.Sprintf("Error when retrieving history data of target %s: %s", content.Target, err.Error()))

			if err := sendMessage(conn, NewMessage(ERR, content.Target, "Internal server error!")); err != nil {
				return fmt.Errorf("%w: %s", FatalWebsocketError, err)
			}

			return err
		}

	} else {
		var err error
		history, err = reader.GetHistoryEntriesFromUntil(string(content.Target), params.From, params.Until)
		if err != nil {
			logger.Error(fmt.Sprintf("Error when retrieving history data of target %s: %s", content.Target, err.Error()))

			if err := sendMessage(conn, NewMessage(ERR, content.Target, "Internal server error!")); err != nil {
				return fmt.Errorf("%w: %s", FatalWebsocketError, err)
			}

			return err
		}
	}

	historyJson, err := json.Marshal(history)

	if err != nil {
		logger.Error(fmt.Sprintf("Error when marshalling history data of target %s: %s", content.Target, err.Error()))

		if err := sendMessage(conn, NewMessage(ERR, content.Target, "Internal server error!")); err != nil {
			return fmt.Errorf("%w: %s", FatalWebsocketError, err)
		}

		return err
	}

	logger.Trace(fmt.Sprintf("Retrieved history of target %s for %s.", content.Target, clientAddress))

	if err = sendMessage(conn, NewMessage(REPLY, content.Target, string(historyJson))); err != nil {
		return fmt.Errorf("%w: %s", FatalWebsocketError, err)
	}

	return nil
}

func sendMessage(conn net.Conn, msg Message) error {
	var err error

	bytes, err := msg.Bytes()
	if err != nil {
		logger.Error(fmt.Sprintf("Encoding %s message for %s was unsuccessful with reason %s.", msg.OpCode, conn.RemoteAddr(), err))
		return err
	}

	err = wsutil.WriteServerText(conn, bytes)
	if err != nil {
		logger.Error(fmt.Sprintf("Sending %s message for %s was unsuccessful with reason %s.", msg.OpCode, conn.RemoteAddr(), err))
		return err
	}

	return nil
}
