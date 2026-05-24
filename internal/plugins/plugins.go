package plugins

import (
	"log"
	"encoding/json"
)

const (
	PluginBuiltin = "builtin"
)

type Plugin interface {
	Type() string
	Execute() (string, error)
}

type Constructor func(rawConfig json.RawMessage) Plugin

type Registry map[string]Plugin
type Constructors map[string]Constructor

var (
	registry     = Registry{}
	constructors = Constructors{}
)

func Register(name string, plugin Plugin) {
	log.Printf("registered plugin: %v", name)
	registry[name] = plugin
}

func RegisterConstructor(
	name string,
	constructor Constructor,
) {
	log.Printf("registered plugin constructor: %v", name)
	constructors[name] = constructor
}

func Get(name string) Plugin {
	return registry[name]
}

func GetConstructor(name string) Constructor {
	return constructors[name]
}
