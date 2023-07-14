package pubsub

// Message is used to model pubsub messages
type Message struct {
	monitor string
	body    string
}

// NewMessage constructs a new message
func NewMessage(message string, monitor string) *Message {
	return &Message{monitor, message}
}

// GetMonitor gives back the name of the monitor the message was flagged with
func (message *Message) GetMonitor() string {
	return message.monitor
}

// GetMessageBody returns the body of the message
func (message *Message) GetMessageBody() string {
	return message.body
}
