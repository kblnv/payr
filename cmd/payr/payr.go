package main

import (
	"fmt"
	"os"
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

	if params.Init {
		runInit(params.ConfigPath)
		return
	}

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

func runInit(configPath string) {
	config := `{
    "server": {
      "address": "127.0.0.1:8080"
    },

    "plugins_dir": "./plugins",

    "handlers": {
      "healthcheck": {
        "plugin": "template",
        "settings": {
          "template": "healthcheck"
        }
      }
    },

    "transports": {
      "telegram": {
        "bot_token": "<bot_token>",
        "chat_id": "<chat_id>"
      }
    },

    "events": [
      {
        "name": "healthcheck",
        "transports": ["telegram"],
        "handler": "healthcheck"
      }
    ]
  }`

	if err := os.WriteFile(configPath, []byte(config), 0644); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create config: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Ready! Run: ./payr")
}
