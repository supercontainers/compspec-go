package plugins

// A Plugin(s)Information interface is an easy way to combine plugins across spaces
// primarily to expose metadata, etc.
type PluginsInformation interface {
	GetPlugins() []PluginInformation
}

type PluginInformation interface {
	GetName() string
	GetType() string
	GetSections() []PluginSection
	GetDescription() string
}

type PluginSection struct {
	Description string
	Name        string
}
