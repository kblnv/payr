package domain

import (
	"encoding/json"

	"payr/internal/repository"
)

type Event struct {
	Name     string
	Plugin   string
	Settings json.RawMessage
}

type Registry struct {
	Events     map[string]Event
	Transports map[string]json.RawMessage
}

type GlobalSettings struct {
	Host    string
	Port    string
	Plugins string
}

func GetGlobalSettings(registryDTO *repository.Registry) *GlobalSettings {
	return &GlobalSettings{
		Host:    registryDTO.Server.Host,
		Port:    registryDTO.Server.Port,
		Plugins: registryDTO.Plugins,
	}
}

func GetRegistry(registryDTO *repository.Registry) *Registry {
	registry := Registry{
		Events:     make(map[string]Event, len(registryDTO.Events)),
		Transports: make(map[string]json.RawMessage, len(registryDTO.Transports)),
	}

	for name, e := range registryDTO.Events {
		registry.Events[name] = Event{
			Name:     e.Name,
			Plugin:   e.Handler.Plugin,
			Settings: e.Handler.Settings,
		}
	}

	for key, t := range registryDTO.Transports {
		registry.Transports[key] = t
	}

	return &registry
}
