package repository

import (
	"encoding/json"
	"os"

	"payr/internal/logger"
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

type Config struct {
	Path   string
	Logger *logger.Logger
}

type Repository struct {
	path string
	log  *logger.Logger
}

func New(config Config) *Repository {
	return &Repository{
		path: config.Path,
		log:  config.Logger,
	}
}

func (c *Repository) GetAll() *Registry {
	bytes, err := os.ReadFile(c.path)
	if err != nil {
		c.log.Fatal("failed to read config file: %v", err)
	}

	var registry Registry

	err = json.Unmarshal(bytes, &registry)
	if err != nil {
		c.log.Fatal("failed to parse config: %v", err)
	}

	return &registry
}
