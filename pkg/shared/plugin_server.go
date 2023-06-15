package shared

import "github.com/Excubitor-Monitoring/Excubitor-Backend/pkg/shared/modules"

type ModuleRPCServer struct {
	Impl ModuleProvider
}

func (s *ModuleRPCServer) GetName(_ []interface{}, response *string) error {
	*response = s.Impl.GetName()
	return nil
}

func (s *ModuleRPCServer) GetVersion(_ []interface{}, response *modules.Version) error {
	*response = s.Impl.GetVersion()
	return nil
}

func (s *ModuleRPCServer) TickFunction(_ []interface{}, response *[]PluginMessage) error {
	*response = s.Impl.TickFunction()
	return nil
}

func (s *ModuleRPCServer) GetComponents(_ []interface{}, response *[]modules.Component) error {
	*response = s.Impl.GetComponents()
	return nil
}
