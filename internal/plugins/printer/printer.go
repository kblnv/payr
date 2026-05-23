package printer

import (
	"payr/internal/plugins"
)

type Printer struct{}

func (p *Printer) Name() string {
	return "printer"
}

func (p *Printer) Type() string {
	return plugins.PluginBuiltin
}

func (p *Printer) Execute() (string, error) {
	return "printer", nil
}

func New() *Printer {
	return &Printer{}
}
