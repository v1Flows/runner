package plugins

import (
	"net/rpc"
)

// Plugin interface that all plugins must implement
type Plugin interface {
	Execute(args map[string]string) (string, error)
	Info() (PluginInfo, error)
}

// PluginInfo holds metadata about the plugin
type PluginInfo struct {
	Name    string
	Version string
	Author  string
}

// PluginRPC is an implementation of net/rpc for Plugin
type PluginRPC struct {
	client *rpc.Client
}

func (p *PluginRPC) Execute(args map[string]string) (string, error) {
	var resp string
	err := p.client.Call("Plugin.Execute", args, &resp)
	return resp, err
}

func (p *PluginRPC) Info() (PluginInfo, error) {
	var resp PluginInfo
	err := p.client.Call("Plugin.Info", new(interface{}), &resp)
	return resp, err
}
