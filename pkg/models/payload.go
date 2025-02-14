package models

type PayloadEndpoint struct {
	Name     string `json:"name"`
	Type     string `json:"type"`
	Endpoint string `json:"endpoint"`
	Version  string `json:"version"`
}
