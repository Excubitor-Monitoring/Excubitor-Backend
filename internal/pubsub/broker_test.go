package pubsub

import (
	"github.com/Excubitor-Monitoring/Excubitor-Backend/internal/logging"
	"github.com/stretchr/testify/assert"
	"os"
	"sync"
	"testing"
)

func TestMain(m *testing.M) {
	err := logging.SetDefaultLogger("CONSOLE")
	if err != nil {
		panic(err)
	}
	code := m.Run()
	os.Exit(code)
}

func TestNewBroker(t *testing.T) {
	broker := NewBroker()

	assert.IsType(t, &Broker{}, broker)
}

func TestAddSubscriber(t *testing.T) {
	broker := NewBroker()
	sub := broker.AddSubscriber()

	assert.IsType(t, &Subscriber{}, sub)
	assert.Empty(t, sub.GetMonitors())
}

func TestSubscribe(t *testing.T) {
	broker := NewBroker()
	sub := broker.AddSubscriber()

	broker.Subscribe(sub, "Some Monitor")
	assert.Contains(t, sub.GetMonitors(), "Some Monitor")
}

func TestUnsubscribe(t *testing.T) {
	broker := NewBroker()
	sub := broker.AddSubscriber()

	broker.Subscribe(sub, "Some Monitor")
	assert.Contains(t, sub.GetMonitors(), "Some Monitor")

	broker.Unsubscribe(sub, "Some Monitor")
	assert.NotContains(t, sub.GetMonitors(), "Some Monitor")
	assert.Empty(t, sub.GetMonitors())

	broker.Subscribe(sub, "Some other monitor")
	assert.Contains(t, sub.GetMonitors(), "Some other monitor")

	broker.Unsubscribe(sub, "Some monitor")
	assert.Contains(t, sub.GetMonitors(), "Some other monitor")
}

func TestPublish(t *testing.T) {
	wg := sync.WaitGroup{}
	wg1 := sync.WaitGroup{}
	wg.Add(1)
	wg1.Add(1)

	broker := NewBroker()
	sub := broker.AddSubscriber()
	sub1 := broker.AddSubscriber()

	broker.Subscribe(sub, "Monitor")
	broker.Subscribe(sub1, "Monitor")

	go sub.Listen(func(message *Message) {
		assert.Equal(t, message.GetMonitor(), "Monitor")
		assert.Equal(t, message.GetMessageBody(), "Test Message!")
		wg.Done()
	})

	go sub1.Listen(func(message *Message) {
		assert.Equal(t, message.GetMonitor(), "Monitor")
		assert.Equal(t, message.GetMessageBody(), "Test Message!")
		wg1.Done()
	})

	broker.Publish("Monitor", "Test Message!")

	wg.Wait()
	wg1.Wait()

}
