package protocol

import "github.com/AlertFlow/runner/pkg/models"

type Request struct {
	Action string                 `json:"action"`
	Data   map[string]interface{} `json:"data"`
}

type Response struct {
	Success bool                   `json:"success"`
	Data    map[string]interface{} `json:"data"`
	Plugin  models.Plugin          `json:"plugin,omitempty"`
	Error   string                 `json:"error,omitempty"`
}
