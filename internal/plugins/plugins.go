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

func Register(plugin Plugin) {
	pluginsRegistry[plugin.Name()] = plugin
}

func GetAll() PluginsRegistry {
	return pluginsRegistry
}
