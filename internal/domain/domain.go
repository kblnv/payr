package domain

import (
	"payr/internal/repository"
)

type Plugin struct {
	Type string
	Name string
}

type Transport struct {
	Sender    string
	ChannelId string
}

type Event struct {
	Transports []string
	Plugin     Plugin
}

type Registry struct {
	Events     map[string]Event
	Transports map[string]Transport
}

func MapRegistry(registryDTO *repository.Registry) (*Registry, error) {
	registry := Registry{
		Events:     make(map[string]Event, len(registryDTO.Events)),
		Transports: make(map[string]Transport, len(registryDTO.Transports)),
	}

	for _, e := range registryDTO.Events {
		registry.Events[e.Name] = Event{
			Transports: e.Transports,
			Plugin: Plugin{
				Type: e.Plugin.Type,
				Name: e.Plugin.Name,
			},
		}
	}

	for _, t := range registryDTO.Transports {
		registry.Transports[t.Name] = Transport{
			Sender:    t.Sender,
			ChannelId: t.ChannelId,
		}
	}

	return &registry, nil
}
