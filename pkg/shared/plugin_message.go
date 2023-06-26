package shared

// PluginMessage is a model of a message a plugin may send to the main application when ticked.
type PluginMessage struct {
	Monitor string
	Body    string
}
