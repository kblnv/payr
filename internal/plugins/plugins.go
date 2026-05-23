package plugins

const (
	PluginBuiltin = "builtin"
)

type Plugin interface {
	Name() string
	Type() string
	Execute() (string, error)
}

type PluginsRegistry map[string]Plugin

var pluginsRegistry = PluginsRegistry{}

func RegisterPlugin(plugin Plugin) {
	pluginsRegistry[plugin.Name()] = plugin
}

func GetPlugins() PluginsRegistry {
	return pluginsRegistry
}
