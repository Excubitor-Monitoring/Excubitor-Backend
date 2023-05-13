package websocket

import (
	"encoding/json"
	"fmt"
	ctx "github.com/Excubitor-Monitoring/Excubitor-Backend/internal/context"
	"github.com/Excubitor-Monitoring/Excubitor-Backend/internal/db"
	"github.com/Excubitor-Monitoring/Excubitor-Backend/internal/logging"
	"github.com/Excubitor-Monitoring/Excubitor-Backend/internal/pubsub"
	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"net"
	"sync"
)

var logger logging.Logger

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
		err = sendMessage(conn, NewMessage(REPLY, TargetAddress(m.GetMonitor()), m.GetMessageBody()))
		if err != nil {
			logger.Error(fmt.Sprintf("Couldn't send message to %s. Aborting connection...", clientAddress))
			return
		}
	})

	for {
		// Receiving message

		msg, op, err := wsutil.ReadClientData(conn)

		if op == ws.OpClose {
			logger.Debug(fmt.Sprintf("Client from %s closes connection.", clientAddress))
			err := conn.Close()
			if err != nil {
				logger.Error(fmt.Sprintf("Couldn't close connection from %s after client closed websocket!", clientAddress))
				return
			}
		}

		if err != nil {
			logger.Warn(fmt.Sprintf("Can't receive message from %s! Aborting connection...", clientAddress))
			return
		}

		logger.Trace(fmt.Sprintf("Received message from %s: %s", clientAddress, string(msg)))

		// Decoding message

		content := &Message{}
		if err = json.Unmarshal(msg, content); err != nil {
			logger.Warn(fmt.Sprintf("Can't decode message from %s! Dropping request...", clientAddress))

			err := sendMessage(conn, NewMessage(ERR, GetEmptyTarget(), "Bad Request!"))
			if err != nil {
				logger.Error(fmt.Sprintf("Sending error message for %s was unsuccessful with reason %s. Forcing connection to close.", clientAddress, err.Error()))
				return
			}

			continue
		}

		switch content.OpCode {
		case GET:
			if err := handleGET(conn, content); err != nil {
				return
			}
		case SUB:
			broker.Subscribe(subscriber, string(content.Target))
			logger.Trace("Client", clientAddress, "subscribed to monitor", content.Target)
		case UNSUB:
			broker.Unsubscribe(subscriber, string(content.Target))
			logger.Trace("Client", clientAddress, "unsubscribed from monitor", content.Target)
		case HIST:
			if err := handleHIST(conn, content); err != nil {
				return
			}
		case REPLY:
			err = sendMessage(conn, NewMessage(ERR, content.Target, "Clients may not send messages of the type REPLY!"))
			if err != nil {
				logger.Warn(fmt.Sprintf("Sending error message to %s was unsuccessful! Aborting connection...", clientAddress))
				return
			}
		case ERR:
			err = sendMessage(conn, NewMessage(ERR, content.Target, "Clients may not send messages of the type ERROR"))
			if err != nil {
				logger.Warn(fmt.Sprintf("Sending error message to %s was unsuccessful! Aborting connection...", clientAddress))
				return
			}
		default:
			err = sendMessage(conn, NewMessage(ERR, content.Target, fmt.Sprintf("Unsupported Operation %s!", content.OpCode)))
			if err != nil {
				logger.Warn(fmt.Sprintf("Sending error message to %s was unsuccessful! Aborting connection...", clientAddress))
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
				logger.Error(fmt.Sprintf("Couldn't send message to %s. Aborting connection...", conn.RemoteAddr()))
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
	reader := db.GetReader()
	history, err := reader.GetHistoryEntriesByTarget(string(content.Target))
	if err != nil {
		logger.Error(fmt.Sprintf("Error when retrieving history data of target %s: %s", content.Target, err.Error()))
		err := sendMessage(conn, NewMessage(ERR, content.Target, "Internal server error!"))
		if err != nil {
			logger.Error(fmt.Sprintf("Sending error message for %s was unsuccessful with reason. Forcing connection to close.", clientAddress))
			return err
		}
	}

	historyJson, err := json.Marshal(history)
	if err != nil {
		logger.Error(fmt.Sprintf("Error when marshalling history data of target %s: %s", content.Target, err.Error()))
		err := sendMessage(conn, NewMessage(ERR, content.Target, "Internal server error!"))
		if err != nil {
			logger.Error(fmt.Sprintf("Sending error message for %s was unsuccessful with reason. Forcing connection to close.", clientAddress))
			return err
		}
	}

	logger.Trace(fmt.Sprintf("Retrieved history of target %s for %s.", content.Target, clientAddress))

	err = sendMessage(conn, NewMessage(REPLY, content.Target, string(historyJson)))
	if err != nil {
		logger.Error(fmt.Sprintf("Couldn't send message to %s. Aborting connection...", clientAddress))
		return err
	}

	return nil
}

func sendMessage(conn net.Conn, msg Message) error {
	var err error

	bytes, err := msg.Bytes()
	if err != nil {
		return err
	}

	err = wsutil.WriteServerText(conn, bytes)
	if err != nil {
		return err
	}

	return nil
}
