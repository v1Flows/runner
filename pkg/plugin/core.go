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

type GRPCServer struct {
	pb.UnimplementedAlertFlowPluginServer
	Impl PluginHandler
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

func (p *Plugin) GRPCServer(broker *plugin.GRPCBroker, s *grpc.Server) error {
	pb.RegisterAlertFlowPluginServer(s, &GRPCServer{Impl: p.Impl})
	return nil
}

func (p *Plugin) GRPCClient(ctx context.Context, broker *plugin.GRPCBroker, c *grpc.ClientConn) (interface{}, error) {
	return &GRPCClient{client: pb.NewAlertFlowPluginClient(c)}, nil
}
