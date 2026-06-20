package server

import (
	"encoding/json"
	"net/http"

	"payr/internal/domain"
	"payr/internal/logger"
	"payr/internal/plugins"
	"payr/internal/transports"

	api "payr/pkg/plugins"
)

type Server struct {
	server            *http.Server
	registry          *domain.Registry
	pluginsManager    *plugins.Manager
	transportsManager *transports.Transports
	log               *logger.Logger
}

type Config struct {
	Logger            *logger.Logger
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
		log:               config.Logger,
	}

	mux.HandleFunc("/event", server.handleEventTrigger)

	return &server
}

func (s *Server) Start() {
	s.log.Info("server is listening on %v", s.server.Addr)

	err := s.server.ListenAndServe()
	if err != nil {
		s.log.Fatal("server error: %v", err)
	}
}

func (s *Server) handleEventTrigger(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	defer r.Body.Close()

	var payload EventTriggerRequestBody
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		s.log.Warn("failed to parse request body: %v", err)
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	s.log.Info("handling event: %v", payload.Event)

	event, ok := s.registry.Events[payload.Event]
	if !ok {
		s.log.Warn("event not found: %v", payload.Event)
		http.Error(w, "event not found", http.StatusNotFound)
		return
	}

	plugin := s.pluginsManager.Get(event.Handler)
	if plugin == nil {
		s.log.Error("plugin not found: %v", event.Handler)
		http.Error(w, "plugin not found", http.StatusInternalServerError)
		return
	}

	ctx := api.Context{
		EventMeta: payload.Meta,
	}

	result, err := plugin.Execute(&ctx)
	if err != nil {
		s.log.Error("plugin execution failed: %v", err)
		http.Error(w, "plugin execution failed", http.StatusInternalServerError)
		return
	}

	s.log.Debug("plugin result: %v", result)

	for _, name := range event.Transports {
		transport := s.transportsManager.Get(name)
		if transport == nil {
			s.log.Error("transport not found: %v", name)
			http.Error(w, "transport not found", http.StatusInternalServerError)
			return
		}

		err := transport.Send(result)
		if err != nil {
			s.log.Error("transport send failed: %v", err)
			http.Error(w, "transport send failed", http.StatusInternalServerError)
			return
		}
	}

	w.Write([]byte("ok"))
}
