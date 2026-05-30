package main

import (
	"bytes"
	"encoding/json"
	"text/template"

	"payr/pkg/plugins"
)

type Config struct {
	Template string `json:"template"`
}

type TemplatePlugin struct {
	template *template.Template
}

func New(rawConfig json.RawMessage) plugins.Plugin {
	var config Config

	err := json.Unmarshal(rawConfig, &config)
	if err != nil {
		panic(err)
	}

	tmpl, err := template.New("message").Parse(config.Template)
	if err != nil {
		panic(err)
	}

	return &TemplatePlugin{
		template: tmpl,
	}
}

func (t *TemplatePlugin) Execute(context *plugins.Context) (string, error) {
	var data map[string]any

	if len(context.EventMeta) > 0 {
		if err := json.Unmarshal(context.EventMeta, &data); err != nil {
			return "", err
		}
	} else {
		data = map[string]any{}
	}

	var buf bytes.Buffer

	if err := t.template.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}
