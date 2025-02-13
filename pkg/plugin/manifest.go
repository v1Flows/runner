package plugin

import "time"

type PluginManifest struct {
	Name        string    `json:"name"`
	Version     string    `json:"version"`
	URL         string    `json:"url"`
	SHA256      string    `json:"sha256"`
	LastUpdated time.Time `json:"last_updated"`
}

type PluginRegistry struct {
	Plugins []PluginManifest `json:"plugins"`
}
