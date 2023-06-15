package modules

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
