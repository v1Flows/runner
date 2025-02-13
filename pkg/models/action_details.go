package models

type ActionDetails struct {
	ID                string      `json:"id"`
	Name              string      `json:"name"`
	Description       string      `json:"description"`
	Version           string      `json:"version"`
	Icon              string      `json:"icon"`
	Type              string      `json:"type"`
	Category          string      `json:"category"`
	Function          interface{} `json:"-"`
	IsHidden          bool        `json:"is_hidden"`
	Params            []Param     `json:"params"`
	CustomName        string      `json:"custom_name"`
	CustomDescription string      `json:"custom_description"`
}

type Param struct {
	Key         string `json:"key"`
	Type        string `json:"type"`
	Required    bool   `json:"required"`
	Default     any    `json:"default"`
	Description string `json:"description"`
}
