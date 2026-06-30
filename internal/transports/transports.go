package transports

import (
	"encoding/json"

	"payr/internal/logger"
)

type Transport interface {
	Send(text string, to string) error
}

type Constructor func(log *logger.Logger, rawConfig json.RawMessage) (Transport, error)

type Registry map[string]Transport
type Constructors map[string]Constructor

type Transports struct {
	registry     Registry
	constructors Constructors
	log          *logger.Logger
}

func (t *Transports) Register(name string, transport Transport) {
	t.log.Info("registered transport: %v", name)
	t.registry[name] = transport
}

func (t *Transports) Get(name string) Transport {
	return t.registry[name]
}

func (t *Transports) RegisterConstructor(name string, constructor Constructor) {
	t.constructors[name] = constructor
}

func (t *Transports) GetConstructor(name string) Constructor {
	return t.constructors[name]
}

func New(log *logger.Logger) *Transports {
	return &Transports{
		registry:     Registry{},
		constructors: Constructors{},
		log:          log,
	}
}
