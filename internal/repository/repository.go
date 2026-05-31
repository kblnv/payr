package repository

import (
	"encoding/json"
	"os"

	"payr/internal/helpers"
)

type Plugin struct {
	Type     string          `json:"type"`
	Settings json.RawMessage `json:"settings"`
}

type Event struct {
	Name       string   `json:"name"`
	Transports []string `json:"transports"`
	Plugin     string   `json:"plugin"`
}

type Server struct {
	Address string `json:"address"`
}

type Registry struct {
	Server     Server                     `json:"server"`
	PluginsDir string                     `json:"plugins_dir"`
	Plugins    map[string]Plugin          `json:"plugins"`
	Transports map[string]json.RawMessage `json:"transports"`
	Events     []Event                    `json:"events"`
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

func (c *Repository) GetAll() *Registry {
	bytes, err := os.ReadFile(c.settings.Path)
	helpers.Die(err)

	var registry Registry

	err = json.Unmarshal(bytes, &registry)
	helpers.Die(err)

	return &registry
}
