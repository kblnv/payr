package plugins

const (
	PluginBuiltin = "builtin"
)

type Plugin interface {
	Name() string
	Type() string
	Execute() (string, error)
}

type Registry map[string]Plugin

var registry = Registry{}

func Register(plugin Plugin) {
	registry[plugin.Name()] = plugin
}

func Get(name string) Plugin {
	return registry[name]
}
