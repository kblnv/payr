package domain

import (
	"encoding/json"

	"payr/internal/repository"
)

type Plugin struct {
	Type     string
	Settings json.RawMessage
}

type Event struct {
	Transports []string
	Plugin     string
}

type Registry struct {
	Events     map[string]Event
	Transports map[string]json.RawMessage
	Plugins    map[string]Plugin
}

type GlobalSettings struct {
	ServerAddress string
	PluginsDir    string
}

func GetGlobalSettings(registryDTO *repository.Registry) *GlobalSettings {
	return &GlobalSettings{
		ServerAddress: registryDTO.Server.Address,
		PluginsDir:    registryDTO.PluginsDir,
	}
}

func GetRegistry(registryDTO *repository.Registry) *Registry {
	registry := Registry{
		Events:     make(map[string]Event, len(registryDTO.Events)),
		Transports: make(map[string]json.RawMessage, len(registryDTO.Transports)),
		Plugins:    make(map[string]Plugin, len(registryDTO.Plugins)),
	}

	for _, e := range registryDTO.Events {
		registry.Events[e.Name] = Event{
			Transports: e.Transports,
			Plugin:     e.Plugin,
		}
	}

	for key, t := range registryDTO.Transports {
		registry.Transports[key] = t
	}

	for key, p := range registryDTO.Plugins {
		registry.Plugins[key] = Plugin{
			Type:     p.Type,
			Settings: p.Settings,
		}
	}

	return &registry
}
