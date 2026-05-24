package main

import (
	"encoding/json"
	"log"
	"net/http"

	"payr/internal/cmd"
	"payr/internal/domain"
	"payr/internal/helpers"
	"payr/internal/repository"

	"payr/internal/transports"
	_ "payr/internal/transports/exports"

	"payr/internal/plugins"
	_ "payr/internal/plugins/exports"
)

const (
	DEFAULT_CONFIG_PATH    = "./registry.json"
	DEFAULT_SERVER_ADDRESS = "127.0.0.1:8080"
)

type EventRequestBody struct {
	Event string `json:"event"`
}

func handleEvent(event domain.Event) error {
	if event.Plugin.Type != "builtin" {
		helpers.Todo("plugins are in development")
	}

	plugin := plugins.Get(event.Plugin.Name)
	result, err := plugin.Execute()

	if err != nil {
		return err
	}

	for _, name := range event.Transports {
		transport := transports.Get(name)
		err := transport.Send(result)

		if err != nil {
			return err
		}
	}

	return nil
}

func main() {
	cmd := cmd.New(cmd.Params{
		ConfigPath:    DEFAULT_CONFIG_PATH,
		ServerAddress: DEFAULT_SERVER_ADDRESS,
	})

	params := cmd.Parse()

	repository := repository.New(repository.Settings{
		Path: params.ConfigPath,
	})

	registryDTO, err := repository.GetRegistry()
	helpers.Die(err)

	registry, err := domain.MapRegistry(registryDTO)
	helpers.Die(err)

	for _, event := range registry.Events {
		constructor := plugins.GetConstructor(event.Plugin.Name)
		plugin := constructor(event.Plugin.Settings)

		plugins.Register(event.Plugin.Name, plugin)
	}

	for name, config := range registry.Transports {
		constructor := transports.GetConstructor(name)
		transport := constructor(config)

		transports.Register(name, transport)
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/event", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		defer r.Body.Close()

		var payload EventRequestBody
		err := json.NewDecoder(r.Body).Decode(&payload)
		helpers.Die(err)

		log.Printf("handling event: %v...", payload.Event)
		event := registry.Events[payload.Event]

		err = handleEvent(event)
		helpers.Die(err)

		w.Write([]byte("ok"))
	})

	server := &http.Server{
		Addr:    params.ServerAddress,
		Handler: mux,
	}

	log.Printf("server is listening on %v...", params.ServerAddress)

	err = server.ListenAndServe()
	helpers.Die(err)
}
