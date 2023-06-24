package shared

import (
	"github.com/Excubitor-Monitoring/Excubitor-Backend/pkg/shared/modules"
)

// ModuleRPCServer is the server of the plugin connection.
// Its methods are run on the plugin side and called through ModuleRPC in the main program.
type ModuleRPCServer struct {
	Impl ModuleProvider
}

func (s *ModuleRPCServer) GetName(_ interface{}, response *string) error {
	*response = s.Impl.GetName()
	return nil
}

func (s *ModuleRPCServer) GetVersion(_ interface{}, response *modules.Version) error {
	*response = s.Impl.GetVersion()
	return nil
}

func (s *ModuleRPCServer) TickFunction(_ interface{}, response *[]PluginMessage) error {
	*response = s.Impl.TickFunction()
	return nil
}

func (s *ModuleRPCServer) GetComponents(_ interface{}, response *[]modules.Component) error {
	*response = s.Impl.GetComponents()
	return nil
}

func (s *ModuleRPCServer) GetComponentFile(args PathArgs, response *[]byte) error {
	*response = s.Impl.GetComponentFile(args.Path)
	return nil
}
