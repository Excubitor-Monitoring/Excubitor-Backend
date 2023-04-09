package pubsub

import (
	"fmt"
	"github.com/google/uuid"
	"sync"
)

type Listener func(*Message)

// Subscriber can listen to different monitors on a broker. Its messages channel will be updated whenever a new message is published with the associated broker.
type Subscriber struct {
	id       string
	messages chan *Message
	monitors map[string]bool
	active   bool
	lock     sync.RWMutex
}

func newSubscriber() (string, *Subscriber) {
	id := uuid.New().String()

	return id, &Subscriber{
		id:       id,
		messages: make(chan *Message),
		monitors: map[string]bool{},
		active:   true,
	}
}

func (subscriber *Subscriber) addMonitor(monitor string) {
	subscriber.lock.RLock()
	defer subscriber.lock.RUnlock()

	subscriber.monitors[monitor] = true
}

func (subscriber *Subscriber) removeMonitor(monitor string) {
	subscriber.lock.RLock()
	defer subscriber.lock.RUnlock()

	delete(subscriber.monitors, monitor)
}

// GetMonitors returns a slice of all monitors the Subscriber is subscribed to.
func (subscriber *Subscriber) GetMonitors() []string {
	subscriber.lock.RLock()
	defer subscriber.lock.RUnlock()

	monitors := []string{}
	for monitor := range subscriber.monitors {
		monitors = append(monitors, monitor)
	}

	return monitors
}

func (subscriber *Subscriber) signal(message *Message) {
	subscriber.lock.RLock()
	defer subscriber.lock.RUnlock()

	if subscriber.active {
		subscriber.messages <- message
	}
}

// Listen listens for messages on the messages channel and calls a Listener function with the message as an argument.
func (subscriber *Subscriber) Listen(listener Listener) {
	for {
		if message, ok := <-subscriber.messages; ok {
			logger.Trace(fmt.Sprintf("Subscriber %s received message from %s: %s", subscriber.id, message.GetMonitor(), message.GetMessageBody()))
			listener(message)
		}
	}
}

// Destruct destructs a Subscriber
func (subscriber *Subscriber) Destruct() {
	subscriber.lock.RLock()
	defer subscriber.lock.RUnlock()

	subscriber.active = false
	close(subscriber.messages)
}
