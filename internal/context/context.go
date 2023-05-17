package ctx

import (
	"github.com/Excubitor-Monitoring/Excubitor-Backend/internal/logging"
	"github.com/Excubitor-Monitoring/Excubitor-Backend/internal/pubsub"
	"sync"
)

var singletonOnce sync.Once

type Module struct {
	Name         string `json:"name"`
	tickFunction func()
}

func NewModule(name string, tickFunction func()) *Module {
	return &Module{name, tickFunction}
}

type Context struct {
	broker  *pubsub.Broker
	modules map[string]*Module
	logger  logging.Logger
	lock    sync.RWMutex
}

var context *Context

func GetContext() *Context {
	if context == nil {
		singletonOnce.Do(func() {
			context = &Context{
				modules: map[string]*Module{},
			}
		})
	}

	return context
}

func (ctx *Context) RegisterModule(module *Module) {
	ctx.lock.RLock()
	defer ctx.lock.RUnlock()

	ctx.modules[module.Name] = module

	startClock()
}

func (ctx *Context) GetModules() []Module {
	ctx.lock.RLock()
	defer ctx.lock.RUnlock()

	var modules []Module
	for _, module := range ctx.modules {
		modules = append(modules, *module)
	}

	return modules
}

func (ctx *Context) RegisterBroker(broker *pubsub.Broker) {
	ctx.lock.RLock()
	defer ctx.lock.RUnlock()

	ctx.broker = broker
}

func (ctx *Context) GetBroker() *pubsub.Broker {
	ctx.lock.RLock()
	defer ctx.lock.RUnlock()

	return ctx.broker
}
