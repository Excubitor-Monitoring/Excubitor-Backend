package shared

import (
	"fmt"
	"github.com/Excubitor-Monitoring/Excubitor-Backend/internal/logging"
	"github.com/Excubitor-Monitoring/Excubitor-Backend/pkg/shared/modules"
	"net/rpc"
)

// ModuleRPC is the client of the plugin connection.
// Its methods are run in the main program to call the respective methods on the plugin side.
type ModuleRPC struct {
	client *rpc.Client
}

func (rpc *ModuleRPC) GetName() string {
	var response string
	if err := rpc.client.Call("Plugin.GetName", new(interface{}), &response); err != nil {
		logging.GetLogger().Error(fmt.Sprintf("Error when calling 'Plugin.GetName' over RPC: %s", err))
		panic(err)
	}

	return response
}

func (rpc *ModuleRPC) GetVersion() modules.Version {
	var response modules.Version
	if err := rpc.client.Call("Plugin.GetVersion", new(interface{}), &response); err != nil {
		logging.GetLogger().Error(fmt.Sprintf("Error when calling 'Plugin.GetVersion' over RPC: %s", err))
		panic(err)
	}

	return response
}

func (rpc *ModuleRPC) GetComponents() []modules.Component {
	var response []modules.Component
	if err := rpc.client.Call("Plugin.GetComponents", new(interface{}), &response); err != nil {
		logging.GetLogger().Error(fmt.Sprintf("Error when calling 'Plugin.GetComponents' over RPC: %s", err))
		panic(err)
	}

	return response
}

func (rpc *ModuleRPC) GetComponentFile(path string) []byte {
	var response []byte
	if err := rpc.client.Call("Plugin.GetComponentFile", PathArgs{path}, &response); err != nil {
		logging.GetLogger().Error(fmt.Sprintf("Error when calling 'Plugin.GetComponentFile' over RPC: %s", err))
		panic(err)
	}

	return response
}

func (rpc *ModuleRPC) TickFunction() []PluginMessage {
	var response []PluginMessage
	if err := rpc.client.Call("Plugin.TickFunction", new(interface{}), &response); err != nil {
		logging.GetLogger().Error(fmt.Sprintf("Error when calling 'Plugin.TickFunction' over RPC: %s", err))
		panic(err)
	}

	return response
}
