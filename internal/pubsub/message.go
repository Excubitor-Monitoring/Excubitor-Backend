package pubsub

type Message struct {
	monitor string
	body    string
}

func NewMessage(message string, monitor string) *Message {
	return &Message{monitor, message}
}

func (message *Message) GetMonitor() string {
	return message.monitor
}

func (message *Message) GetMessageBody() string {
	return message.body
}
