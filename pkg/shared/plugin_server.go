package shared

import "github.com/Excubitor-Monitoring/Excubitor-Backend/pkg/shared/modules"

type ModuleRPCServer struct {
	Impl ModuleProvider
}

func (s *ModuleRPCServer) GetName() string {
	return s.Impl.GetName()
}

func (s *ModuleRPCServer) GetVersion() modules.Version {
	return s.Impl.GetVersion()
}

func (s *ModuleRPCServer) TickFunction() []PluginMessage {
	return s.Impl.TickFunction()
}

func (s *ModuleRPCServer) GetComponents() []modules.Component {
	return s.Impl.GetComponents()
}
