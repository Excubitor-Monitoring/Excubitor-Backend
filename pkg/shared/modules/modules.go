package modules

type Module struct {
	Name         string      `json:"name"`
	Version      string      `json:"version"`
	Components   []Component `json:"components"`
	TickFunction func()      `json:"-"`
}

func NewModule(name string, version Version, components []Component, tickFunction func()) *Module {
	return &Module{
		Name:         name,
		Version:      version.string(),
		Components:   components,
		TickFunction: tickFunction,
	}
}
