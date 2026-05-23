package repository

import (
	"encoding/json"
	"os"
)

type Plugin struct {
	Type string `json:"type"`
	Name string `json:"name"`
}

type Transport struct {
	Name      string `json:"name"`
	Sender    string `json:"sender"`
	ChannelId string `json:"channel_id"`
}

type Event struct {
	Name       string   `json:"name"`
	Transports []string `json:"transports"`
	Plugin     Plugin   `json:"plugin"`
}

type Registry struct {
	Events     []Event     `json:"events"`
	Transports []Transport `json:"transports"`
}

type Settings struct {
	Path string
}

type Repository struct {
	settings Settings
}

func New(settings Settings) *Repository {
	return &Repository{settings: settings}
}

func (c *Repository) GetRegistry() (*Registry, error) {
	bytes, err := os.ReadFile(c.settings.Path)

	if err != nil {
		return nil, err
	}

	var registry Registry

	if err := json.Unmarshal(bytes, &registry); err != nil {
		return nil, err
	}

	return &registry, nil
}
