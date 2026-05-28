package main

import (
	"encoding/json"

	"payr/pkg/plugins"
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
	
	if (err != nil) {
		panic(err)
	}

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
