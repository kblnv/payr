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
	DEFAULT_CONFIG_PATH = "./registry.json"
)

func main() {
	cmd := cmd.New(cmd.InParams{
		ConfigPath: DEFAULT_CONFIG_PATH,
	})
	params := cmd.Parse()

	repository := repository.New(repository.Settings{
		Path: params.ConfigPath,
	})
	config := repository.GetAll()

	registry := domain.GetRegistry(config)
	globalSettings := domain.GetGlobalSettings(config)

	pluginsManager := plugins.New()
	pluginsManager.LoadAll(globalSettings.PluginsDir)

	for _, event := range registry.Events {
		config := registry.Plugins[event.Plugin]

		constructor := pluginsManager.GetConstructor(config.Type)
		instance := constructor(config.Settings)

		pluginsManager.Register(event.Plugin, instance)
	}

	transportsManager := transports.New()

	for name, config := range registry.Transports {
		constructor := transports.GetConstructor(name)
		transport := constructor(config)

		transportsManager.Register(name, transport)
	}

	server := server.New(server.Config{
		Address:           globalSettings.ServerAddress,
		Registry:          registry,
		PluginsManager:    pluginsManager,
		TransportsManager: transportsManager,
	})

	server.Start()
}
