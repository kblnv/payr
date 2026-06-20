package plugins

import (
	"encoding/json"
	"os"
	"path/filepath"
	"plugin"
	"strings"

	"payr/internal/logger"
	"payr/pkg/plugins"
)

type Registry map[string]plugins.Plugin
type Constructors map[string]plugins.Constructor

type Manager struct {
	registry     Registry
	constructors Constructors
	log          *logger.Logger
}

func (m *Manager) Get(name string) plugins.Plugin {
	return m.registry[name]
}

func (m *Manager) GetConstructor(name string) plugins.Constructor {
	return m.constructors[name]
}

func (m *Manager) Register(name string, plugin plugins.Plugin) {
	m.log.Info("registered handler: %v", name)
	m.registry[name] = plugin
}

func (m *Manager) RegisterConstructor(name string, constructor plugins.Constructor) {
	m.constructors[name] = constructor
}

func (m *Manager) LoadAll(path string) {
	files, err := os.ReadDir(path)
	if err != nil {
		m.log.Fatal("failed to read plugins directory: %v", err)
	}

	for _, file := range files {
		fileName := file.Name()

		if !strings.HasSuffix(fileName, ".so") {
			continue
		}

		fullPath := filepath.Join(path, fileName)

		pkg, err := plugin.Open(fullPath)
		if err != nil {
			m.log.Fatal("failed to open plugin: %v", err)
		}

		m.log.Debug("loaded plugin: %v", fullPath)

		sym, err := pkg.Lookup("New")
		if err != nil {
			m.log.Fatal("failed to find plugin constructor: %v", err)
		}

		constructor, ok := sym.(func(json.RawMessage) (plugins.Plugin, error))
		if !ok {
			m.log.Fatal("invalid plugin constructor in %s", fullPath)
		}

		name := strings.TrimSuffix(fileName, filepath.Ext(fileName))

		m.RegisterConstructor(name, constructor)
	}
}

func New(log *logger.Logger) *Manager {
	return &Manager{
		registry:     Registry{},
		constructors: Constructors{},
		log:          log,
	}
}
