package transports

import (
	"encoding/json"
	"log"
)

type Transport interface {
	Send(text string) error
}

type Constructor func(rawConfig json.RawMessage) Transport

type Registry map[string]Transport
type Constructors map[string]Constructor

type Transports struct {
	registry Registry
}

func (t *Transports) Register(name string, transport Transport) {
	log.Printf("registered transport: %v", name)
	t.registry[name] = transport
}

func (t *Transports) Get(name string) Transport {
	return t.registry[name]
}

func New() *Transports {
	return &Transports{
		registry: Registry{},
	}
}

var constructors = Constructors{}

func RegisterConstructor(
	name string,
	constructor Constructor,
) {
	log.Printf("registered transport constructor: %v", name)
	constructors[name] = constructor
}

func GetConstructor(name string) Constructor {
	return constructors[name]
}
