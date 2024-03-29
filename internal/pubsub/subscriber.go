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
	wg       sync.WaitGroup
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
	subscriber.lock.Lock()
	defer subscriber.lock.Unlock()

	subscriber.monitors[monitor] = true
}

func (subscriber *Subscriber) removeMonitor(monitor string) {
	subscriber.lock.Lock()
	defer subscriber.lock.Unlock()

	delete(subscriber.monitors, monitor)
}

// GetMonitors returns a slice of all monitors the Subscriber is subscribed to.
func (subscriber *Subscriber) GetMonitors() []string {
	subscriber.lock.RLock()
	defer subscriber.lock.RUnlock()

	var monitors []string
	for monitor := range subscriber.monitors {
		monitors = append(monitors, monitor)
	}

	return monitors
}

func (subscriber *Subscriber) signal(message *Message) {
	subscriber.lock.RLock()
	defer subscriber.lock.RUnlock()

	if subscriber.active {
		subscriber.wg.Add(1)
		defer subscriber.wg.Done()

		subscriber.messages <- message
	}
}

// Listen listens for messages on the messages channel and calls a Listener function with the message as an argument.
func (subscriber *Subscriber) Listen(listener Listener) {
	for {
		if !subscriber.active {
			break
		}

		if message, ok := <-subscriber.messages; ok {
			logger.Trace(fmt.Sprintf("Subscriber %s received message from %s.", subscriber.id, message.GetMonitor()))
			listener(message)
		}
	}
}

// Destruct destructs a Subscriber
func (subscriber *Subscriber) Destruct() {
	subscriber.lock.RLock()
	defer subscriber.lock.RUnlock()

	logger.Trace(fmt.Sprintf("Destructing subscriber %s", subscriber.id))

	subscriber.active = false

	go func() {
		subscriber.wg.Wait()
		close(subscriber.messages)
	}()
}
