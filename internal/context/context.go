package ctx

import (
	"github.com/Excubitor-Monitoring/Excubitor-Backend/internal/logging"
	"github.com/Excubitor-Monitoring/Excubitor-Backend/internal/pubsub"
	"github.com/Excubitor-Monitoring/Excubitor-Backend/pkg/shared/modules"
	"sync"
)

var singletonOnce sync.Once

type Context struct {
	broker  *pubsub.Broker
	modules map[string]*modules.Module
	logger  logging.Logger
	lock    sync.RWMutex
}

var context *Context

func GetContext() *Context {
	if context == nil {
		singletonOnce.Do(func() {
			context = &Context{
				modules: map[string]*modules.Module{},
			}
		})
	}

	return context
}

func (ctx *Context) RegisterModule(module *modules.Module) {
	ctx.lock.RLock()
	defer ctx.lock.RUnlock()

	ctx.modules[module.Name] = module

	startClock()
}

func (ctx *Context) GetModules() []modules.Module {
	ctx.lock.RLock()
	defer ctx.lock.RUnlock()

	var modules []modules.Module
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
