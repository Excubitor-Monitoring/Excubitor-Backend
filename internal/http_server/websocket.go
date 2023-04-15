package http_server

import (
	"encoding/json"
	"fmt"
	ctx "github.com/Excubitor-Monitoring/Excubitor-Backend/internal/context"
	"github.com/Excubitor-Monitoring/Excubitor-Backend/internal/pubsub"
	"golang.org/x/net/websocket"
)

type OpCode int

const (
	GET OpCode = iota
	SUB
	UNSUB
	HIST
	REPLY
	ERR
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

func handleWebsocket(ws *websocket.Conn) {
	var err error
	clientAddress := ws.Request().RemoteAddr

	broker := ctx.GetContext().GetBroker()
	subscriber := broker.AddSubscriber()

	subscriber.Listen(func(m *pubsub.Message) {
		err = sendMessage(ws, newMessage(REPLY, TargetAddress(m.GetMonitor()), m.GetMessageBody()))
		if err != nil {
			logger.Error(fmt.Sprintf("Couldn't send message to %s. Aborting connection...", clientAddress))

			err := ws.Close()
			if err != nil {
				logger.Error(fmt.Sprintf("Connection to %s couldn't be closed gracefully. Exiting connection...", clientAddress))
			}

			return
		}
	})

	for {
		var request []byte

		if err = websocket.Message.Receive(ws, &request); err != nil {
			logger.Warn(fmt.Sprintf("Can't receive message from %s! Aborting connection...", clientAddress))

			err = sendMessage(ws, newMessage(ERR, GetEmptyTarget(), "Undecipherable Message!"))
			if err != nil {
				logger.Error(fmt.Sprintf("Sending error message to %s was unsuccessful with reason %s. Forcing connection to close.", clientAddress, err.Error()))

				err := ws.Close()
				if err != nil {
					logger.Error(fmt.Sprintf("WebSocket connection to %s couldn't be closed.", clientAddress))
				}

				return
			}

			continue
		}

		logger.Trace(fmt.Sprintf("Received message from %s: %s", clientAddress, string(request)))

		content := &message{}
		if err = json.Unmarshal(request, content); err != nil {
			logger.Warn(fmt.Sprintf("Can't decode message from %s! Dropping request...", clientAddress))

			err := sendMessage(ws, newMessage(ERR, GetEmptyTarget(), "Bad Request!"))
			if err != nil {
				logger.Error(fmt.Sprintf("Sending error message for %s was unsuccessful with reason %s. Forcing connection to close.", clientAddress, err.Error()))
				return
			}

			continue
		}

		switch content.OpCode {
		case GET:
			err = sendMessage(ws, newMessage(ERR, content.Target, "This feature is not implemented yet!"))
			if err != nil {
				logger.Warn("Sending error message to", clientAddress, "was unsuccessful! Aborting connection...")
				return
			}
			break
		case SUB:
			broker.Subscribe(subscriber, string(content.Target))
			logger.Trace("Client", clientAddress, "subscribed to monitor", content.Target)
			break
		case UNSUB:
			broker.Unsubscribe(subscriber, string(content.Target))
			logger.Trace("Client", clientAddress, "unsubscribed from monitor", content.Target)
			break
		case HIST:
			err = sendMessage(ws, newMessage(ERR, content.Target, "This feature is not implemented yet!"))
			if err != nil {
				logger.Warn("Sending error message to", clientAddress, "was unsuccessful! Aborting connection...")
				return
			}
			break
		case REPLY:
			err = sendMessage(ws, newMessage(ERR, content.Target, "Clients may not send messages of the type REPLY!"))
			break
		case ERR:
			err = sendMessage(ws, newMessage(ERR, content.Target, "Clients may not send messages of the type ERROR"))
			break
		}
	}
}

func sendMessage(ws *websocket.Conn, msg message) error {
	var err error

	bytes, err := msg.bytes()
	if err != nil {
		return err
	}

	_, err = ws.Write(bytes)
	if err != nil {
		return err
	}

	return nil
}
