package plugins

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"plugin"
	"strings"

	"payr/internal/helpers"
	"payr/pkg/plugins"
)

type Registry map[string]plugins.Plugin
type Constructors map[string]plugins.Constructor

type Manager struct {
	registry     Registry
	constructors Constructors
}

func (m *Manager) Get(name string) plugins.Plugin {
	return m.registry[name]
}

func (m *Manager) GetConstructor(name string) plugins.Constructor {
	return m.constructors[name]
}

func (m *Manager) Register(name string, plugin plugins.Plugin) {
	log.Printf("registered plugin: %v", name)
	m.registry[name] = plugin
}

func (m *Manager) RegisterConstructor(name string, constructor plugins.Constructor) {
	log.Printf("registered plugin constructor: %v", name)
	m.constructors[name] = constructor
}

func (m *Manager) LoadAll(path string) {
	files, err := os.ReadDir(path)
	helpers.Die(err)

	for _, file := range files {
		fileName := file.Name()

		if !strings.HasSuffix(fileName, ".so") {
			continue
		}

		fullPath := filepath.Join(path, fileName)

		pkg, err := plugin.Open(fullPath)
		helpers.Die(err)

		log.Println("loaded plugin:", fullPath)

		sym, err := pkg.Lookup("New")
		helpers.Die(err)

		constructor, ok := sym.(func(json.RawMessage) plugins.Plugin)
		if !ok {
			log.Fatalf("invalid plugin constructor in %s", fullPath)
		}

		name := strings.TrimSuffix(fileName, filepath.Ext(fileName))

		m.RegisterConstructor(name, constructor)
	}
}

func New() *Manager {
	return &Manager{
		registry:     Registry{},
		constructors: Constructors{},
	}
}
