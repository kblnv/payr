package server

import (
	"encoding/json"
	"log"
	"net/http"

	"payr/internal/domain"
	"payr/internal/helpers"
	"payr/internal/plugins"
	"payr/internal/transports"

	api "payr/pkg/plugins"
)

type Server struct {
	server            *http.Server
	registry          *domain.Registry
	pluginsManager    *plugins.Manager
	transportsManager *transports.Transports
}

type Config struct {
	Address           string
	Registry          *domain.Registry
	PluginsManager    *plugins.Manager
	TransportsManager *transports.Transports
}

type EventTriggerRequestBody struct {
	Event string          `json:"event"`
	Meta  json.RawMessage `json:"meta"`
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

	plugin := s.pluginsManager.Get(event.Plugin)

	ctx := api.Context{
		EventMeta: payload.Meta,
	}

	result, err := plugin.Execute(&ctx)
	helpers.Die(err)

	log.Println("plugin result:", result)

	for _, name := range event.Transports {
		transport := s.transportsManager.Get(name)
		err := transport.Send(result)

		helpers.Die(err)
	}

	w.Write([]byte("ok"))
}
