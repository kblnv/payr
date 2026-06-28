package cmd

import (
	"flag"
)

type InParams struct {
	ConfigPath string
}

type OutParams struct {
	ConfigPath string
	Init       bool
}

type Cmd struct {
	defaultConfigPath string
}

func New(defaultParams InParams) *Cmd {
	return &Cmd{
		defaultConfigPath: defaultParams.ConfigPath,
	}
}

func (c *Cmd) Parse() OutParams {
	configPath := flag.String(
		"config",
		c.defaultConfigPath,
		"config path",
	)
	init := flag.Bool(
		"init",
		false,
		"initialize new config",
	)
	flag.Parse()

	return OutParams{
		ConfigPath: *configPath,
		Init:       *init,
	}
}
