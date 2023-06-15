package shared

import (
	"embed"
	"github.com/Excubitor-Monitoring/Excubitor-Backend/pkg/shared/modules"
	"github.com/hashicorp/go-plugin"
	"net/rpc"
)

type ModuleProvider interface {
	GetName() string
	GetVersion() modules.Version
	TickFunction() []PluginMessage
	GetComponents() []modules.Component
	GetComponentFiles() embed.FS
}

type ModulePlugin struct {
	Impl ModuleProvider
}

func (p *ModulePlugin) Server(*plugin.MuxBroker) (interface{}, error) {
	return &ModuleRPCServer{Impl: p.Impl}, nil
}

func (*ModulePlugin) Client(_ *plugin.MuxBroker, c *rpc.Client) (interface{}, error) {
	return &ModuleRPC{c}, nil
}
