package cmd

import (
	"flag"
)

type InParams struct {
	ConfigPath    string
	ServerAddress string
}

type OutParams struct {
	ConfigPath    string
	ServerAddress string
	PluginsDir    string
}

type Cmd struct {
	defaultConfigPath    string
	defaultServerAddress string
}

func New(defaultParams InParams) *Cmd {
	return &Cmd{
		defaultConfigPath:    defaultParams.ConfigPath,
		defaultServerAddress: defaultParams.ServerAddress,
	}
}

func (c *Cmd) Parse() OutParams {
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

	pluginsDir := flag.String(
		"plugins",
		"",
		"plugins directory",
	)

	flag.Parse()

	return OutParams{
		ConfigPath:    *configPath,
		ServerAddress: *serverAddress,
		PluginsDir:    *pluginsDir,
	}
}
