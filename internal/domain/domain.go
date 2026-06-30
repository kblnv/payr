package domain

import (
	"encoding/json"

	"payr/internal/repository"
)

type Handler struct {
	Plugin   string
	Settings json.RawMessage
}

type Event struct {
	Transports []string
	Handler    string
}

type Registry struct {
	Events     map[string]Event
	Transports map[string]json.RawMessage
	Handlers   map[string]Handler
}

type GlobalSettings struct {
	Host       string
	Port       string
	PluginsDir string
}

func GetGlobalSettings(registryDTO *repository.Registry) *GlobalSettings {
	return &GlobalSettings{
		Host:       registryDTO.Server.Host,
		Port:       registryDTO.Server.Port,
		PluginsDir: registryDTO.PluginsDir,
	}
}

func GetRegistry(registryDTO *repository.Registry) *Registry {
	registry := Registry{
		Events:     make(map[string]Event, len(registryDTO.Events)),
		Transports: make(map[string]json.RawMessage, len(registryDTO.Transports)),
		Handlers:   make(map[string]Handler, len(registryDTO.Handlers)),
	}

	for _, e := range registryDTO.Events {
		registry.Events[e.Name] = Event{
			Transports: e.Transports,
			Handler:    e.Handler,
		}
	}

	for key, t := range registryDTO.Transports {
		registry.Transports[key] = t
	}

	for key, p := range registryDTO.Handlers {
		registry.Handlers[key] = Handler{
			Plugin:   p.Plugin,
			Settings: p.Settings,
		}
	}

	return &registry
}
