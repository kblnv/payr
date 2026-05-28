package plugins

import "encoding/json"

type Plugin interface {
	Execute() (string, error)
}

type Constructor func(rawConfig json.RawMessage) Plugin
