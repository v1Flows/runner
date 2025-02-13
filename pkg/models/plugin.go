package models

type Plugin struct {
	Name    string `json:"name"`
	Type    string `json:"type"`
	Version string `json:"version"`
	Author  string `json:"author"`
	Payload PayloadEndpoint
	Action  ActionDetails
}
