package models

type Plugin struct {
	Name    string `json:"name"`
	Type    string `json:"type"`
	Version string `json:"version"`
	Creator string `json:"creator"`
}
