package plugin

import (
	"context"

	pb "github.com/AlertFlow/runner/pkg/plugin/proto"
	"github.com/hashicorp/go-plugin"
	"google.golang.org/grpc"
)

type PluginHandler interface {
	Execute(ctx context.Context, req *pb.PluginRequest) (*pb.PluginResponse, error)
	StreamStatus(req *pb.PluginRequest, stream pb.AlertFlowPlugin_StreamStatusServer) error
}

type Plugin struct {
	Impl PluginHandler
}

// GRPCPlugin implements the go-plugin.GRPCPlugin interface
type GRPCPlugin struct {
	plugin.Plugin
	Impl PluginHandler
}

type GRPCServer struct {
	pb.UnimplementedAlertFlowPluginServer
	Impl PluginHandler
}

// Handshake is a shared configuration between the host and plugins
var Handshake = plugin.HandshakeConfig{
	ProtocolVersion:  1,
	MagicCookieKey:   "ALERTFLOW_PLUGIN",
	MagicCookieValue: "d7e562f1-3d77-4c84-9e87-f4f19e766d0d",
}

func (s *GRPCServer) Execute(ctx context.Context, req *pb.PluginRequest) (*pb.PluginResponse, error) {
	return s.Impl.Execute(ctx, req)
}

func (s *GRPCServer) StreamStatus(req *pb.PluginRequest, stream pb.AlertFlowPlugin_StreamStatusServer) error {
	return s.Impl.StreamStatus(req, stream)
}

type GRPCClient struct {
	client pb.AlertFlowPluginClient
}

func (c *GRPCClient) Execute(ctx context.Context, req *pb.PluginRequest) (*pb.PluginResponse, error) {
	return c.client.Execute(ctx, req)
}

func (c *GRPCClient) StreamStatus(req *pb.PluginRequest, server pb.AlertFlowPlugin_StreamStatusServer) error {
	stream, err := c.client.StreamStatus(context.Background(), req)
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

// GRPCServer registers the plugin server implementation
func (p *GRPCPlugin) GRPCServer(_ *plugin.GRPCBroker, s *grpc.Server) error {
	pb.RegisterAlertFlowPluginServer(s, &GRPCServer{Impl: p.Impl})
	return nil
}

// GRPCClient creates a new client implementation
func (p *GRPCPlugin) GRPCClient(_ context.Context, _ *plugin.GRPCBroker, c *grpc.ClientConn) (interface{}, error) {
	return &GRPCClient{client: pb.NewAlertFlowPluginClient(c)}, nil
}
