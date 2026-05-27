package server

import (
	"encoding/json"
	"log"
	"net/http"

	"payr/internal/domain"
	"payr/internal/helpers"
	"payr/internal/plugins"
	"payr/internal/transports"
)

type Server struct {
	server            *http.Server
	registry          *domain.Registry
	pluginsManager    *plugins.Plugins
	transportsManager *transports.Transports
}

type Config struct {
	Address           string
	Registry          *domain.Registry
	PluginsManager    *plugins.Plugins
	TransportsManager *transports.Transports
}

type EventTriggerRequestBody struct {
	Event string `json:"event"`
}

func New(config Config) *Server {
	mux := http.NewServeMux()

	server := Server{
		server: &http.Server{
			Addr:    config.Address,
			Handler: mux,
		},
		registry:          config.Registry,
		pluginsManager:    config.PluginsManager,
		transportsManager: config.TransportsManager,
	}

	mux.HandleFunc("/event", server.handleEventTrigger)

	return &server
}

func (s *Server) Start() {
	log.Printf("server is listening on %v...", s.server.Addr)

	err := s.server.ListenAndServe()
	helpers.Die(err)
}

func (s *Server) handleEventTrigger(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	defer r.Body.Close()

	var payload EventTriggerRequestBody
	err := json.NewDecoder(r.Body).Decode(&payload)
	helpers.Die(err)

	log.Printf("handling event: %v...", payload.Event)
	event := s.registry.Events[payload.Event]

	if event.Plugin.Type != "builtin" {
		helpers.Todo("plugins are in development")
	}

	plugin := s.pluginsManager.Get(event.Plugin.Name)
	result, err := plugin.Execute()
	helpers.Die(err)

	for _, name := range event.Transports {
		transport := s.transportsManager.Get(name)
		err := transport.Send(result)

		helpers.Die(err)
	}

	w.Write([]byte("ok"))
}
