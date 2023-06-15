package shared

import (
	"github.com/Excubitor-Monitoring/Excubitor-Backend/pkg/shared/modules"
	"net/rpc"
)

type ModuleRPC struct {
	client *rpc.Client
}

func (rpc *ModuleRPC) GetName() string {
	var response string
	if err := rpc.client.Call("Plugin.GetName", new(interface{}), &response); err != nil {
		panic(err) // TODO: Better error handling
	}

	return response
}

func (rpc *ModuleRPC) GetVersion() modules.Version {
	var response modules.Version
	if err := rpc.client.Call("Plugin.GetVersion", new(interface{}), &response); err != nil {
		panic(err) // TODO: Better error handling
	}

	return response
}

func (rpc *ModuleRPC) GetComponents() []modules.Component {
	var response []modules.Component
	if err := rpc.client.Call("Plugin.GetName", new(interface{}), &response); err != nil {
		panic(err) // TODO: Better error handling
	}

	return response
}

func (rpc *ModuleRPC) TickFunction() []PluginMessage {
	var response []PluginMessage
	if err := rpc.client.Call("Plugin.GetName", new(interface{}), &response); err != nil {
		panic(err) // TODO: Better error handling
	}

	return response
}
