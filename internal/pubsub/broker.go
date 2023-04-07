package pubsub

import (
	"fmt"
	"github.com/Excubitor-Monitoring/Excubitor-Backend/internal/logging"
	"sync"
)

var logger logging.Logger

type Subscribers map[string]*Subscriber

type Broker struct {
	subscribers Subscribers
	monitors    map[string]Subscribers
	lock        sync.RWMutex
}

func NewBroker() *Broker {
	return &Broker{
		subscribers: Subscribers{},
		monitors:    map[string]Subscribers{},
	}
}

func (broker *Broker) AddSubscriber() *Subscriber {
	broker.lock.Lock()
	defer broker.lock.Unlock()

	id, subscriber := newSubscriber()
	logger.Trace(fmt.Sprintf("Adding new subscriber with id %s.", id))

	broker.subscribers[id] = subscriber
	return subscriber
}

func (broker *Broker) Subscribe(subscriber *Subscriber, monitor string) {
	broker.lock.Lock()
	defer broker.lock.Unlock()

	logger.Debug(fmt.Sprintf("Subscribing %s to %s.", subscriber.id, monitor))

	if broker.monitors[monitor] == nil {
		broker.monitors[monitor] = Subscribers{}
	}

	subscriber.addMonitor(monitor)
	broker.monitors[monitor][subscriber.id] = subscriber
}

func (broker *Broker) Unsubscribe(subscriber *Subscriber, monitor string) {
	broker.lock.RLock()
	defer broker.lock.RUnlock()

	logger.Debug(fmt.Sprintf("Unsubscribing %s from monitor %s.", monitor, subscriber.id))

	delete(broker.monitors[monitor], subscriber.id)
	subscriber.removeMonitor(monitor)
}

func (broker *Broker) Publish(monitor string, message string) {
	broker.lock.RLock()
	subscribers := broker.monitors[monitor]
	broker.lock.RUnlock()

	logger.Trace(fmt.Sprintf("Publishing message %s on monitor %s.", monitor, message))

	for _, subscriber := range subscribers {
		m := NewMessage(message, monitor)
		if !subscriber.active {
			return
		}

		go (func(s *Subscriber) {
			s.signal(m)
		})(subscriber)
	}
}

func init() {
	var err error
	logger, err = logging.GetLogger()

	if err != nil {
		panic(err)
	}
}
