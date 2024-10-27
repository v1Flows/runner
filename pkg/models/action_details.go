package models

import (
	"encoding/json"
)

type ActionDetails struct {
	Name              string          `json:"name"`
	Description       string          `json:"description"`
	Icon              string          `json:"icon"`
	Type              string          `json:"type"`
	Function          interface{}     `json:"-"`
	Params            json.RawMessage `json:"params"`
	CustomName        string          `json:"custom_name"`
	CustomDescription string          `json:"custom_description"`
}

type Param struct {
	Key         string `json:"key"`
	Type        string `json:"type"`
	Required    bool   `json:"required"`
	Default     any    `json:"default"`
	Description string `json:"description"`
}
