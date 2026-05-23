package main

import (
	"encoding/json"
	"log"
	"net/http"

	"payr/internal/cmd"
	"payr/internal/domain"
	"payr/internal/helpers"
	"payr/internal/repository"
	"payr/internal/plugins"

	_ "payr/internal/plugins/export"
)

const (
	DEFAULT_CONFIG_PATH    = "./registry.json"
	DEFAULT_SERVER_ADDRESS = "127.0.0.1:8080"
)

func handleEvent(pluginsRegistry plugins.PluginsRegistry, event domain.Event) (string, error) {
	if (event.Plugin.Type != "builtin" || event.Plugin.Name != "printer") {
		helpers.Todo("plugins are in development")
	}

	plugin := pluginsRegistry[event.Plugin.Name]
	result, err := plugin.Execute()

	return result, err
}

type EventRequestBody struct {
	Event string `json:"event"`
}

func main() {
	pluginsRegistry := plugins.GetPlugins()

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

		event := registry.Events[payload.Event]

		result, err := handleEvent(pluginsRegistry, event)
		helpers.Die(err)

		w.Write([]byte(result))
	})

	server := &http.Server{
		Addr:    params.ServerAddress,
		Handler: mux,
	}

	log.Println("server is running...")

	err = server.ListenAndServe()
	helpers.Die(err)
}
