package main

import (
	"payr/internal/cmd"
	"payr/internal/domain"
	"payr/internal/logger"
	"payr/internal/plugins"
	"payr/internal/repository"
	"payr/internal/server"
	"payr/internal/transports"

	"payr/internal/transports/telegram"
)

const (
	DEFAULT_CONFIG_PATH = "./registry.json"
)

func main() {
	log := logger.New("main")

	cmd := cmd.New(cmd.InParams{
		ConfigPath: DEFAULT_CONFIG_PATH,
	})
	params := cmd.Parse()

	repository := repository.New(repository.Config{
		Path:   params.ConfigPath,
		Logger: log,
	})
	config := repository.GetAll()

	registry := domain.GetRegistry(config)
	globalSettings := domain.GetGlobalSettings(config)

	pluginsManager := plugins.New(logger.New("plugins"))
	pluginsManager.LoadAll(globalSettings.PluginsDir)

	for _, event := range registry.Events {
		config := registry.Handlers[event.Handler]

		constructor := pluginsManager.GetConstructor(config.Plugin)
		instance, err := constructor(config.Settings)
		if err != nil {
			log.Fatal("failed to create plugin %s: %v", event.Handler, err)
		}

		pluginsManager.Register(event.Handler, instance)
	}

	transportsManager := transports.New(logger.New("transports"))
	transportsManager.RegisterConstructor("telegram", telegram.New)

	for name, config := range registry.Transports {
		constructor := transportsManager.GetConstructor(name)
		transport, err := constructor(log, config)
		if err != nil {
			log.Fatal("failed to create transport %s: %v", name, err)
		}

		transportsManager.Register(name, transport)
	}

	server := server.New(server.Config{
		Logger:            logger.New("server"),
		Address:           globalSettings.ServerAddress,
		Registry:          registry,
		PluginsManager:    pluginsManager,
		TransportsManager: transportsManager,
	})

	server.Start()
}
