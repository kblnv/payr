package plugins

import (
	"encoding/json"
	"os"
	"path/filepath"
	"plugin"
	"strings"

	"payr/internal/helpers"
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
	m.log.Info("registered plugin: %v", name)
	m.registry[name] = plugin
}

func (m *Manager) RegisterConstructor(name string, constructor plugins.Constructor) {
	m.log.Info("registered plugin constructor: %v", name)
	m.constructors[name] = constructor
}

func (m *Manager) LoadAll(path string) {
	files, err := os.ReadDir(path)
	helpers.Must(err)

	for _, file := range files {
		fileName := file.Name()

		if !strings.HasSuffix(fileName, ".so") {
			continue
		}

		fullPath := filepath.Join(path, fileName)

		pkg, err := plugin.Open(fullPath)
		helpers.Must(err)

		m.log.Debug("loaded plugin: %v", fullPath)

		sym, err := pkg.Lookup("New")
		helpers.Must(err)

		constructor, ok := sym.(func(json.RawMessage) plugins.Plugin)
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
