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
		content, err := newMessage(REPLY, TargetAddress(m.GetMonitor()), m.GetMessageBody()).bytes()
		if err != nil {
			return
		}

		_, err = ws.Write(content)
		if err != nil {
			return
		}
	})

	for {
		var request []byte

		if err = websocket.Message.Receive(ws, &request); err != nil {
			logger.Warn(fmt.Sprintf("Can't receive message from %s! Aborting connection...", clientAddress))

			errMessage, err := newMessage(ERR, GetEmptyTarget(), "Undecipherable Message!").bytes()
			if err != nil {
				logger.Error(fmt.Sprintf("Creating error message to %s was unsuccessful. Forcing connection to close.", clientAddress))
				return
			}

			_, err = ws.Write(errMessage)
			if err != nil {
				logger.Error(fmt.Sprintf("Sending error message to %s was unsuccessful. Forcing connection to close.", clientAddress))
				return
			}

			continue
		}

		logger.Trace(fmt.Sprintf("Received message from %s: %s", clientAddress, string(request)))

		content := &message{}
		if err = json.Unmarshal(request, content); err != nil {
			logger.Warn(fmt.Sprintf("Can't decode message from %s! Dropping request...", clientAddress))

			errMessage, err := newMessage(ERR, GetEmptyTarget(), "Bad Request!").bytes()
			if err != nil {
				logger.Error(fmt.Sprintf("Creating error message for %s was unsuccessful. Forcing connection to close.", clientAddress))
				return
			}

			_, err = ws.Write(errMessage)
			if err != nil {
				logger.Error(fmt.Sprintf("Sending error message to %s was unsuccessful. Forcing connection to close.", clientAddress))
				return
			}

			continue
		}

		switch content.OpCode {
		case GET:
			break
		case SUB:
			broker.Subscribe(subscriber, string(content.Target))
			break
		case UNSUB:
			broker.Unsubscribe(subscriber, string(content.Target))
			break
		case HIST:
			break
		case REPLY:
			break
		case ERR:
			break
		}

		/*if err = websocket.Message.Send(ws, request); err != nil {
			logger.Warn(fmt.Sprintf("Can't send message to %s! Aborting connection...", clientAddress))

			errMessage, err := newMessage(ERR, GetEmptyTarget(), "").bytes()
			if err != nil {
				logger.Error(fmt.Sprintf("Creating error message to %s was unsuccessful. Forcing connection to close.", clientAddress))
				return
			}

			_, err = ws.Write(errMessage)
			if err != nil {
				logger.Error(fmt.Sprintf("Sending error message to %s was unsuccessful. Forcing connection to close.", clientAddress))
				return
			}
		}*/
	}
}
