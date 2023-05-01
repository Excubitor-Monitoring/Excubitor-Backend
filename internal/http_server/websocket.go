package http_server

import (
	"encoding/json"
	"fmt"
	ctx "github.com/Excubitor-Monitoring/Excubitor-Backend/internal/context"
	"github.com/Excubitor-Monitoring/Excubitor-Backend/internal/pubsub"
	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"net"
	"sync"
)

type OpCode string

const (
	GET   OpCode = "GET"
	SUB   OpCode = "SUB"
	UNSUB OpCode = "UNSUB"
	HIST  OpCode = "HIST"
	REPLY OpCode = "REPLY"
	ERR   OpCode = "ERR"
)

type TargetAddress string

func GetEmptyTarget() TargetAddress {
	return ""
}

type message struct {
	OpCode OpCode        `json:"op"`
	Target TargetAddress `json:"target"`
	Value  string        `json:"value"`
}

func newMessage(opcode OpCode, target TargetAddress, value string) message {
	return message{opcode, target, value}
}

func (msg message) bytes() ([]byte, error) {
	jsonData, err := json.Marshal(msg)
	if err != nil {
		return nil, err
	}

	return jsonData, nil
}

func handleWebsocket(conn net.Conn) {
	var err error
	clientAddress := conn.RemoteAddr()

	defer func(conn net.Conn) {
		logger.Debug(fmt.Sprintf("Closing connection from %s", clientAddress))

		err := conn.Close()
		if err != nil {
			logger.Error(fmt.Sprintf("Couldn't close connection from %s", clientAddress))
		}
	}(conn)

	broker := ctx.GetContext().GetBroker()
	subscriber := broker.AddSubscriber()
	defer func(subscriber *pubsub.Subscriber) {
		logger.Trace(fmt.Sprintf("Destructing subscriber associated with connection from %s", clientAddress))
		subscriber.Destruct()
	}(subscriber)

	go subscriber.Listen(func(m *pubsub.Message) {
		logger.Trace(fmt.Sprintf("Sending message from %s to connection from %s", m.GetMonitor(), clientAddress))
		err = sendMessage(conn, newMessage(REPLY, TargetAddress(m.GetMonitor()), m.GetMessageBody()))
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

		content := &message{}
		if err = json.Unmarshal(msg, content); err != nil {
			logger.Warn(fmt.Sprintf("Can't decode message from %s! Dropping request...", clientAddress))

			err := sendMessage(conn, newMessage(ERR, GetEmptyTarget(), "Bad Request!"))
			if err != nil {
				logger.Error(fmt.Sprintf("Sending error message for %s was unsuccessful with reason %s. Forcing connection to close.", clientAddress, err.Error()))
				return
			}

			continue
		}

		switch content.OpCode {
		case GET:
			temporarySubscriber := broker.AddSubscriber()

			var receiveOnce sync.Once

			logger.Trace(fmt.Sprintf("Added temporary subscriber to fulfill GET request from %s on monitor %s.", clientAddress, content.Target))

			go temporarySubscriber.Listen(func(m *pubsub.Message) {
				receiveOnce.Do(func() {
					logger.Trace(fmt.Sprintf("Sending single message from %s to connection from %s", m.GetMonitor(), clientAddress))
					broker.Unsubscribe(temporarySubscriber, string(content.Target))
					defer temporarySubscriber.Destruct()

					err = sendMessage(conn, newMessage(REPLY, TargetAddress(m.GetMonitor()), m.GetMessageBody()))
					if err != nil {
						logger.Error(fmt.Sprintf("Couldn't send message to %s. Aborting connection...", clientAddress))
						return
					}
				})
			})

			broker.Subscribe(temporarySubscriber, string(content.Target))
		case SUB:
			broker.Subscribe(subscriber, string(content.Target))
			logger.Trace("Client", clientAddress, "subscribed to monitor", content.Target)
		case UNSUB:
			broker.Unsubscribe(subscriber, string(content.Target))
			logger.Trace("Client", clientAddress, "unsubscribed from monitor", content.Target)
		case HIST:
			err = sendMessage(conn, newMessage(ERR, content.Target, "This feature is not implemented yet!"))
			if err != nil {
				logger.Warn(fmt.Sprintf("Sending error message to %s was unsuccessful! Aborting connection...", clientAddress))
				return
			}
		case REPLY:
			err = sendMessage(conn, newMessage(ERR, content.Target, "Clients may not send messages of the type REPLY!"))
			if err != nil {
				logger.Warn(fmt.Sprintf("Sending error message to %s was unsuccessful! Aborting connection...", clientAddress))
				return
			}
		case ERR:
			err = sendMessage(conn, newMessage(ERR, content.Target, "Clients may not send messages of the type ERROR"))
			if err != nil {
				logger.Warn(fmt.Sprintf("Sending error message to %s was unsuccessful! Aborting connection...", clientAddress))
				return
			}
		default:
			err = sendMessage(conn, newMessage(ERR, content.Target, fmt.Sprintf("Unsupported Operation %s!", content.OpCode)))
			if err != nil {
				logger.Warn(fmt.Sprintf("Sending error message to %s was unsuccessful! Aborting connection...", clientAddress))
				return
			}
		}
	}
}

func sendMessage(conn net.Conn, msg message) error {
	var err error

	bytes, err := msg.bytes()
	if err != nil {
		return err
	}

	err = wsutil.WriteServerText(conn, bytes)
	if err != nil {
		return err
	}

	return nil
}
