// filepath: /Users/Justin.Neubert/projects/v1flows/alertflow/runner/pkg/plugins/plugin.go
package plugins

import (
	"net/rpc"

	"github.com/AlertFlow/runner/config"
	"github.com/hashicorp/go-plugin"
	"github.com/v1Flows/alertFlow/services/backend/pkg/models"
)

// Plugin interface that all plugins must implement
type Plugin interface {
	ExecuteTask(request ExecuteTaskRequest) (Response, error)
	HandlePayload(request PayloadHandlerRequest) (Response, error)
	Info() (models.Plugins, error)
}

// PluginRPC is an implementation of net/rpc for Plugin
type PluginRPC struct {
	Client *rpc.Client
}

type ExecuteTaskRequest struct {
	Args      map[string]string
	Config    config.Config
	Flow      models.Flows
	Execution models.Executions
	Step      models.ExecutionSteps
	Payload   models.Payloads
}

type PayloadHandlerRequest struct {
	Config config.Config
	Body   []byte
}

type Response struct {
	Data    map[string]interface{}
	Success bool
	Error   string
}

func (p *PluginRPC) ExecuteTask(request ExecuteTaskRequest) (Response, error) {
	var resp Response
	err := p.Client.Call("Plugin.ExecuteTask", request, &resp)
	return resp, err
}

func (p *PluginRPC) HandlePayload(request PayloadHandlerRequest) (Response, error) {
	var resp Response
	err := p.Client.Call("Plugin.HandlePayload", request, &resp)
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

func (s *PluginRPCServer) ExecuteTask(request ExecuteTaskRequest, resp *Response) error {
	result, err := s.Impl.ExecuteTask(request)
	*resp = result
	return err
}
func (s *PluginRPCServer) HandlePayload(request PayloadHandlerRequest, resp *Response) error {
	result, err := s.Impl.HandlePayload(request)
	*resp = result
	return err
}

func (s *PluginRPCServer) Info(args interface{}, resp *models.Plugins) error {
	result, err := s.Impl.Info()
	*resp = result
	return err
}
