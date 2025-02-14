// filepath: /Users/Justin.Neubert/projects/v1flows/alertflow/runner/pkg/plugins/plugin.go
package plugins

import (
	"net/rpc"

	"github.com/AlertFlow/runner/pkg/models"
	"github.com/hashicorp/go-plugin"
)

// Plugin interface that all plugins must implement
type Plugin interface {
	Execute(args map[string]string) (string, error)
	Info() (models.Plugin, error)
}

// PluginRPC is an implementation of net/rpc for Plugin
type PluginRPC struct {
	Client *rpc.Client
}

func (p *PluginRPC) Execute(args map[string]string) (string, error) {
	var resp string
	err := p.Client.Call("Plugin.Execute", args, &resp)
	return resp, err
}

func (p *PluginRPC) Info() (models.Plugin, error) {
	var resp models.Plugin
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

func (s *PluginRPCServer) Execute(args map[string]string, resp *string) error {
	result, err := s.Impl.Execute(args)
	*resp = result
	return err
}

func (s *PluginRPCServer) Info(args interface{}, resp *models.Plugin) error {
	result, err := s.Impl.Info()
	*resp = result
	return err
}
