package main

import (
	"payr/internal/cmd"
	"payr/internal/domain"
	"payr/internal/plugins"
	"payr/internal/repository"
	"payr/internal/server"
	"payr/internal/transports"

	_ "payr/internal/transports/exports"
)

const (
	DEFAULT_CONFIG_PATH    = "./registry.json"
	DEFAULT_SERVER_ADDRESS = "127.0.0.1:8080"
	DEFAULT_PLUGINS_DIR    = "./plugins"
)

func main() {
	cmd := cmd.New(cmd.InParams{
		ConfigPath:    DEFAULT_CONFIG_PATH,
		ServerAddress: DEFAULT_SERVER_ADDRESS,
		PluginsDir:    DEFAULT_PLUGINS_DIR,
	})

	params := cmd.Parse()

	pluginsManager := plugins.New()
	pluginsManager.LoadAll(params.PluginsDir)

	repository := repository.New(repository.Settings{
		Path: params.ConfigPath,
	})

	registryDTO := repository.GetRegistry()
	registry := domain.MapRegistry(registryDTO)

	for _, event := range registry.Events {
		constructor := pluginsManager.GetConstructor(event.Plugin.Name)
		plugin := constructor(event.Plugin.Settings)
		
		pluginsManager.Register(event.Plugin.Name, plugin)
	}
	
	transportsManager := transports.New()

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
