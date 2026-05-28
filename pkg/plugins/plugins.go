package plugins

import "encoding/json"

const (
	PluginBuiltin  = "builtin"
	PluginExternal = "external"
)

type Plugin interface {
	Type() string
	Execute() (string, error)
}

type Constructor func(rawConfig json.RawMessage) Plugin
