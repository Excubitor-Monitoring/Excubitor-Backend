package pubsub

import (
	"fmt"
	"github.com/google/uuid"
	"sync"
)

type Listener func(*Message)

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

func (subscriber *Subscriber) Listen(listener Listener) {
	for {
		if message, ok := <-subscriber.messages; ok {
			logger.Trace(fmt.Sprintf("Subscriber %s received message from %s: %s", subscriber.id, message.GetMonitor(), message.GetMessageBody()))
			listener(message)
		}
	}
}

func (subscriber *Subscriber) Destruct() {
	subscriber.lock.RLock()
	defer subscriber.lock.RUnlock()

	subscriber.active = false
	close(subscriber.messages)
}
