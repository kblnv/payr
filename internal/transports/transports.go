package transports

type Transport interface {
	Name() string
	Send(text string) error
}

type TransportsRegistry map[string]Transport

var transportsRegistry = TransportsRegistry{}

func Register(transport Transport) {
	transportsRegistry[transport.Name()] = transport
}

func GetAll() TransportsRegistry {
	return transportsRegistry
}
