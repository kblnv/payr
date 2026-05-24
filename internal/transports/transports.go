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

var (
	registry     = Registry{}
	constructors = Constructors{}
)

func Register(name string, transport Transport) {
	log.Printf("registered transport: %v", name)
	registry[name] = transport
}

func RegisterConstructor(
	name string,
	constructor Constructor,
) {
	log.Printf("registered transport contructor: %v", name)
	constructors[name] = constructor
}

func Get(name string) Transport {
	return registry[name]
}

func GetConstructor(name string) Constructor {
	return constructors[name]
}