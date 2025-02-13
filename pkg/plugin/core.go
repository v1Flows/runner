package plugin

import (
	"context"

	"github.com/AlertFlow/runner/pkg/models"
	"github.com/hashicorp/go-plugin"
	"google.golang.org/grpc"
)

// Handler defines the interface that plugins must implement
type Handler interface {
	Details() *models.Plugin
	Execute(ctx context.Context, req *Request) (*Response, error)
	StreamUpdates(req *Request, updates Plugin_StreamUpdatesServer) error
}

// GRPCPlugin implements plugin.GRPCPlugin
type GRPCPlugin struct {
	plugin.Plugin
	Impl Handler
}

// GRPCServer implements the gRPC server interface
type GRPCServer struct {
	UnimplementedPluginServer
	Impl Handler
}

// Execute implements the RPC method for plugin execution
func (s *GRPCServer) Execute(ctx context.Context, req *Request) (*Response, error) {
	return s.Impl.Execute(ctx, req)
}

// StreamUpdates implements the RPC method for status updates
func (s *GRPCServer) StreamUpdates(req *Request, stream Plugin_StreamUpdatesServer) error {
	return s.Impl.StreamUpdates(req, stream)
}

// GRPCClient implements the client interface
type GRPCClient struct {
	client PluginClient
}

func (c *GRPCClient) Execute(ctx context.Context, req *Request) (*Response, error) {
	return c.client.Execute(ctx, req)
}

func (c *GRPCClient) StreamUpdates(req *Request, server UpdateServer) error {
	stream, err := c.client.StreamUpdates(context.Background(), req)
	if err != nil {
		return err
	}

	for {
		update, err := stream.Recv()
		if err != nil {
			return err
		}
		if err := server.Send(update); err != nil {
			return err
		}
	}
}

func (p *GRPCPlugin) GRPCServer(broker *plugin.GRPCBroker, s *grpc.Server) error {
	RegisterPluginServer(s, &GRPCServer{Impl: p.Impl})
	return nil
}

func (p *GRPCPlugin) GRPCClient(ctx context.Context, broker *plugin.GRPCBroker, c *grpc.ClientConn) (interface{}, error) {
	return &GRPCClient{client: NewPluginClient(c)}, nil
}

// Handshake is used to verify plugin compatibility
var Handshake = plugin.HandshakeConfig{
	ProtocolVersion:  1,
	MagicCookieKey:   "ALERTFLOW_PLUGIN",
	MagicCookieValue: "alertflow_plugin_v1",
}
