package domain

import (
	"encoding/json"

	"payr/internal/repository"
)

type Plugin struct {
	Name     string
	Settings json.RawMessage
}

type Event struct {
	Transports []string
	Plugin     Plugin
}

type Registry struct {
	Events     map[string]Event
	Transports map[string]json.RawMessage
}

func MapRegistry(registryDTO *repository.Registry) *Registry {
	registry := Registry{
		Events:     make(map[string]Event, len(registryDTO.Events)),
		Transports: make(map[string]json.RawMessage, len(registryDTO.Transports)),
	}

	for _, e := range registryDTO.Events {
		registry.Events[e.Name] = Event{
			Transports: e.Transports,
			Plugin: Plugin{
				Name:     e.Plugin.Name,
				Settings: e.Plugin.Settings,
			},
		}
	}

	for _, t := range registryDTO.Transports {
		registry.Transports[t.Name] = t.Settings
	}

	return &registry
}
