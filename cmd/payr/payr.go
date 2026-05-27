package main

import (
	"payr/internal/cmd"
	"payr/internal/domain"
	"payr/internal/repository"
	"payr/internal/server"

	"payr/internal/transports"
	_ "payr/internal/transports/exports"

	"payr/internal/plugins"
	_ "payr/internal/plugins/exports"
)

const (
	DEFAULT_CONFIG_PATH    = "./registry.json"
	DEFAULT_SERVER_ADDRESS = "127.0.0.1:8080"
)

func main() {
	cmd := cmd.New(cmd.Params{
		ConfigPath:    DEFAULT_CONFIG_PATH,
		ServerAddress: DEFAULT_SERVER_ADDRESS,
	})

	params := cmd.Parse()

	repository := repository.New(repository.Settings{
		Path: params.ConfigPath,
	})

	registryDTO := repository.GetRegistry()
	registry := domain.MapRegistry(registryDTO)

	pluginsManager := plugins.New()
	transportsManager := transports.New()

	for _, event := range registry.Events {
		constructor := plugins.GetConstructor(event.Plugin.Name)
		plugin := constructor(event.Plugin.Settings)

		pluginsManager.Register(event.Plugin.Name, plugin)
	}

	for name, config := range registry.Transports {
		constructor := transports.GetConstructor(name)
		transport := constructor(config)

		transportsManager.Register(name, transport)
	}

	server := server.New(server.Config{
		Address:           params.ServerAddress,
		Registry:          registry,
		PluginsManager:    pluginsManager,
		TransportsManager: transportsManager,
	})

	server.Start()
}
