package plugins

import (
	"encoding/json"
	"log"
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

type Plugins struct {
	registry Registry
}

func (p *Plugins) Get(name string) Plugin {
	return p.registry[name]
}

func (p *Plugins) Register(name string, plugin Plugin) {
	log.Printf("registered plugin: %v", name)
	p.registry[name] = plugin
}

func New() *Plugins {
	return &Plugins{
		registry: Registry{},
	}
}

var constructors = Constructors{}

func RegisterConstructor(
	name string,
	constructor Constructor,
) {
	log.Printf("registered plugin constructor: %v", name)
	constructors[name] = constructor
}

func GetConstructor(name string) Constructor {
	return constructors[name]
}
