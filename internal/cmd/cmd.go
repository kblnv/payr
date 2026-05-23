package cmd

import (
	"flag"
)

type Params struct {
	ConfigPath    string
	ServerAddress string
}

type Cmd struct {
	defaultConfigPath    string
	defaultServerAddress string
}

func New(defaultParams Params) *Cmd {
	return &Cmd{
		defaultConfigPath:    defaultParams.ConfigPath,
		defaultServerAddress: defaultParams.ServerAddress,
	}
}

func (c *Cmd) Parse() Params {
	configPath := flag.String(
		"config",
		c.defaultConfigPath,
		"config path",
	)

	serverAddress := flag.String(
		"address",
		c.defaultServerAddress,
		"server address",
	)

	flag.Parse()

	return Params{
		ConfigPath:    *configPath,
		ServerAddress: *serverAddress,
	}
}
