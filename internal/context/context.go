package ctx

import (
	"github.com/Excubitor-Monitoring/Excubitor-Backend/internal/logging"
	"github.com/Excubitor-Monitoring/Excubitor-Backend/internal/pubsub"
	"sync"
)

var logger logging.Logger

type Module struct {
	Name string `json:"name"`
}

func NewModule(name string) *Module {
	return &Module{name}
}

type Context struct {
	broker  *pubsub.Broker
	modules map[string]*Module
	lock    sync.RWMutex
}

var context *Context

func GetContext() *Context {
	var once sync.Once

	if context == nil {
		once.Do(func() {

			logger.Debug("Instantiating context!")

			context = &Context{
				pubsub.NewBroker(),
				map[string]*Module{},
				sync.RWMutex{},
			}
		})
	}

	return context
}

func (ctx *Context) RegisterModule(module *Module) {
	ctx.lock.RLock()
	defer ctx.lock.RUnlock()

	ctx.modules[module.Name] = module
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

func (ctx *Context) GetBroker() *pubsub.Broker {
	return ctx.broker
}

func init() {
	var err error
	logger, err = logging.GetLogger()

	if err != nil {
		panic(err)
	}
}
