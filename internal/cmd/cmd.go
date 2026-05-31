package cmd

import (
	"flag"
)

type InParams struct {
	ConfigPath string
}

type OutParams struct {
	ConfigPath string
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
	flag.Parse()

	return OutParams{
		ConfigPath: *configPath,
	}
}
