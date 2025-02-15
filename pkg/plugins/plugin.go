// filepath: /Users/Justin.Neubert/projects/v1flows/alertflow/runner/pkg/plugins/plugin.go
package plugins

import (
	"net/rpc"

	"github.com/hashicorp/go-plugin"
	"github.com/v1Flows/alertFlow/services/backend/pkg/models"
)

// Plugin interface that all plugins must implement
type Plugin interface {
	Execute(request ExecuteRequest) (ExecuteResponse, error)
	Info() (models.Plugins, error)
}

// PluginRPC is an implementation of net/rpc for Plugin
type PluginRPC struct {
	Client *rpc.Client
}

type ExecuteRequest struct {
	Args      map[string]string
	Flow      models.Flows
	Execution models.Executions
	Step      models.ExecutionSteps
	Payload   models.Payloads
}

type ExecuteResponse struct {
	Success bool
	Error   string
}

func (p *PluginRPC) Execute(request ExecuteRequest) (ExecuteResponse, error) {
	var resp ExecuteResponse
	err := p.Client.Call("Plugin.Execute", request, &resp)
	return resp, err
}

func (p *PluginRPC) Info() (models.Plugins, error) {
	var resp models.Plugins
	err := p.Client.Call("Plugin.Info", new(interface{}), &resp)
	return resp, err
}

// PluginServer is the implementation of plugin.Plugin interface
type PluginServer struct {
	Impl Plugin
}

func (p *PluginServer) Server(*plugin.MuxBroker) (interface{}, error) {
	return &PluginRPCServer{Impl: p.Impl}, nil
}

func (p *PluginServer) Client(b *plugin.MuxBroker, c *rpc.Client) (interface{}, error) {
	return &PluginRPC{Client: c}, nil
}

// PluginRPCServer is the RPC server for Plugin
type PluginRPCServer struct {
	Impl Plugin
}

func (s *PluginRPCServer) Execute(request ExecuteRequest, resp *ExecuteResponse) error {
	result, err := s.Impl.Execute(request)
	*resp = result
	return err
}

func (s *PluginRPCServer) Info(args interface{}, resp *models.Plugins) error {
	result, err := s.Impl.Info()
	*resp = result
	return err
}
