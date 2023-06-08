package ctx

import (
	"fmt"
	"strings"
)

type Version [3]int

func NewVersion(major int, minor int, sub int) Version {
	return Version{major, minor, sub}
}

func (v Version) string() string {
	return strings.Trim(strings.Join(strings.Fields(fmt.Sprint(v)), "."), "[]")
}

type Component struct {
	TabName string `json:"tab_name"` // Name of the tab in the frontend
	JSFile  string `json:"js"`       // JavaScript file that contains the WebComponent
	Tag     string `json:"tag"`      // Name of the HTML tag to import
}

func NewComponent(tabName string, jsFile string, tag string) Component {
	return Component{
		TabName: tabName,
		JSFile:  jsFile,
		Tag:     tag,
	}
}

type Module struct {
	Name         string      `json:"name"`
	Version      string      `json:"version"`
	Components   []Component `json:"components"`
	tickFunction func()
}

func NewModule(name string, version Version, components []Component, tickFunction func()) *Module {
	return &Module{
		Name:         name,
		Version:      version.string(),
		Components:   components,
		tickFunction: tickFunction,
	}
}
