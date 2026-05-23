package transport

type Transport interface {
	Send(text string) error
}

const (
	TransportTelegram = "telegram"
)
