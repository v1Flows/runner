package models

type ActionDetails struct {
	ID                string  `json:"id"`
	Name              string  `json:"name"`
	Description       string  `json:"description"`
	Icon              string  `json:"icon"`
	Category          string  `json:"category"`
	IsHidden          bool    `json:"is_hidden"`
	Params            []Param `json:"params"`
	CustomName        string  `json:"custom_name"`
	CustomDescription string  `json:"custom_description"`
	Version           string  `json:"version"`
}

type Param struct {
	Key         string `json:"key"`
	Type        string `json:"type"`
	Required    bool   `json:"required"`
	Default     any    `json:"default"`
	Description string `json:"description"`
}
