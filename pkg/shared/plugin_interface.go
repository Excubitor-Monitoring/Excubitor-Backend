package shared

import (
	"github.com/Excubitor-Monitoring/Excubitor-Backend/pkg/shared/modules"
	"github.com/hashicorp/go-plugin"
	"net/rpc"
)

// ModuleProvider serves as the interface plugins need to implement to be used as a module in Excubitor.
type ModuleProvider interface {
	GetName() string
	GetVersion() modules.Version
	TickFunction() []PluginMessage
	GetComponents() []modules.Component
	GetComponentFile(path string) []byte
}

// ModulePlugin is a representation of a plugin
type ModulePlugin struct {
	Impl ModuleProvider
}

// Server provides the server part of the plugin connection
func (p *ModulePlugin) Server(*plugin.MuxBroker) (interface{}, error) {
	return &ModuleRPCServer{Impl: p.Impl}, nil
}

// Client provides the client part of the plugin connection
func (*ModulePlugin) Client(_ *plugin.MuxBroker, c *rpc.Client) (interface{}, error) {
	return &ModuleRPC{c}, nil
}

// Arguments for RPC

// PathArgs serves as argument for the GetComponentFile method in the RPC Server
type PathArgs struct {
	Path string
}
