package text

import (
	"encoding/json"

	"payr/internal/helpers"
	"payr/internal/plugins"
)

type Config struct {
	Text string `json:"text"`
}

type Text struct {
	text string
}

func New(rawConfig json.RawMessage) plugins.Plugin {
	var config Config

	err := json.Unmarshal(rawConfig, &config)
	helpers.Die(err)

	return &Text{
		text: config.Text,
	}
}

func (t *Text) Type() string {
	return plugins.PluginBuiltin
}

func (t *Text) Execute() (string, error) {
	return t.text, nil
}

func init() {
	plugins.RegisterConstructor(
		"text",
		New,
	)
}