package plugin

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os/exec"

	"github.com/AlertFlow/runner/pkg/models"
)

type pluginWrapper struct {
	cmd    *exec.Cmd
	stdin  io.WriteCloser
	stdout io.ReadCloser
}

func (p *pluginWrapper) Execute(ctx context.Context, input map[string]interface{}) (map[string]interface{}, error) {
	// Create message structure
	msg := struct {
		Command string                 `json:"command"`
		Input   map[string]interface{} `json:"input"`
	}{
		Command: "execute",
		Input:   input,
	}

	// Send message to plugin
	if err := json.NewEncoder(p.stdin).Encode(msg); err != nil {
		return nil, fmt.Errorf("failed to send input to plugin: %w", err)
	}

	// Read response
	var response struct {
		Error  string                 `json:"error,omitempty"`
		Output map[string]interface{} `json:"output,omitempty"`
	}

	if err := json.NewDecoder(p.stdout).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to read plugin response: %w", err)
	}

	if response.Error != "" {
		return nil, fmt.Errorf("plugin error: %s", response.Error)
	}

	return response.Output, nil
}

func (p *pluginWrapper) Details() models.Plugin {
	return p.Details()
}
