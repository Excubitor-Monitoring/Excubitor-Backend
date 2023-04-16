package pubsub

import (
	"fmt"
	"github.com/Excubitor-Monitoring/Excubitor-Backend/internal/logging"
	"sync"
)

var logger logging.Logger

type Subscribers map[string]*Subscriber

// Broker is used to interact with the pubsub architecture
type Broker struct {
	subscribers Subscribers
	monitors    map[string]Subscribers
	logger      logging.Logger
	lock        sync.RWMutex
}

// NewBroker constructs a new Broker
func NewBroker() *Broker {
	logger = logging.GetLogger()

	return &Broker{
		subscribers: Subscribers{},
		logger:      logging.GetLogger(),
		monitors:    map[string]Subscribers{},
	}
}

// AddSubscriber adds a new subscriber to the subscriber pool and returns its reference
func (broker *Broker) AddSubscriber() *Subscriber {
	broker.lock.Lock()
	defer broker.lock.Unlock()

	id, subscriber := newSubscriber()
	broker.logger.Trace(fmt.Sprintf("Adding new subscriber with id %s.", id))

	broker.subscribers[id] = subscriber
	return subscriber
}

// Subscribe can add a monitor to a given subscriber
func (broker *Broker) Subscribe(subscriber *Subscriber, monitor string) {
	broker.lock.Lock()
	defer broker.lock.Unlock()

	broker.logger.Debug(fmt.Sprintf("Subscribing %s to %s.", subscriber.id, monitor))

	if broker.monitors[monitor] == nil {
		broker.monitors[monitor] = Subscribers{}
	}

	subscriber.addMonitor(monitor)
	broker.monitors[monitor][subscriber.id] = subscriber
}

// Unsubscribe removes a monitor from a given subscriber
func (broker *Broker) Unsubscribe(subscriber *Subscriber, monitor string) {
	broker.lock.RLock()
	defer broker.lock.RUnlock()

	broker.logger.Debug(fmt.Sprintf("Unsubscribing %s from monitor %s.", monitor, subscriber.id))

	delete(broker.monitors[monitor], subscriber.id)
	subscriber.removeMonitor(monitor)
}

// Publish publishes messages flagged with a given monitor to the pubsub architecture
func (broker *Broker) Publish(monitor string, message string) {
	broker.lock.RLock()
	subscribers := broker.monitors[monitor]
	broker.lock.RUnlock()

	broker.logger.Trace(fmt.Sprintf("Publishing message %s on monitor %s.", message, monitor))

	for _, subscriber := range subscribers {
		m := NewMessage(message, monitor)
		if !subscriber.active {
			continue
		}

		go (func(s *Subscriber) {
			s.signal(m)
		})(subscriber)
	}
}
