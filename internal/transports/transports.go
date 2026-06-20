package transports

import (
	"encoding/json"

	"payr/internal/logger"
)

type Transport interface {
	Send(text string) error
}

type Constructor func(rawConfig json.RawMessage) (Transport, error)

type Registry map[string]Transport
type Constructors map[string]Constructor

type Transports struct {
	registry Registry
	log      *logger.Logger
}

func (t *Transports) Register(name string, transport Transport) {
	t.log.Info("registered transport: %v", name)
	t.registry[name] = transport
}

func (t *Transports) Get(name string) Transport {
	return t.registry[name]
}

func New(log *logger.Logger) *Transports {
	return &Transports{
		registry: Registry{},
		log:      log,
	}
}

var constructors = Constructors{}

func RegisterConstructor(
	name string,
	constructor Constructor,
) {
	constructors[name] = constructor
}

func GetConstructor(name string) Constructor {
	return constructors[name]
}
