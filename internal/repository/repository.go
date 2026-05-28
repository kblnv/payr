package repository

import (
	"encoding/json"
	"os"

	"payr/internal/helpers"
)

type Plugin struct {
	Name     string          `json:"name"`
	Settings json.RawMessage `json:"settings"`
}

type Transport struct {
	Name     string          `json:"name"`
	Settings json.RawMessage `json:"settings"`
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

func (c *Repository) GetRegistry() *Registry {
	bytes, err := os.ReadFile(c.settings.Path)
	helpers.Die(err)

	var registry Registry

	err = json.Unmarshal(bytes, &registry)
	helpers.Die(err)

	return &registry
}
