package notification

import (
	"context"
	"os"

	"saastack/core"
	pb "saastack/interfaces/notification/proto"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// NotificationServiceInterface defines the interface for notification service
type NotificationServiceInterface interface {
	pb.NotificationServiceServer
}

// Registry interface for plugin management
type PluginRegistry interface {
	GetPlugin(interfaceName, pluginName string) (interface{}, bool)
}

// NotificationService implements the NotificationServiceServer interface
type NotificationService struct {
	pb.UnimplementedNotificationServiceServer
	registry PluginRegistry
}

// NewNotificationService creates a new instance of NotificationService
func NewNotificationService(registry PluginRegistry) *NotificationService {
	return &NotificationService{
		registry: registry,
	}
}

func init() {
	service := NewNotificationService(core.GlobalRegistry)
	core.GlobalRegistry.RegisterService("notification", service)
}

// RegisterGRPC registers the service with gRPC server
func (s *NotificationService) RegisterGRPC(server *grpc.Server) {
	pb.RegisterNotificationServiceServer(server, s)
}

// RegisterHTTP registers the service with HTTP gateway
func (s *NotificationService) RegisterHTTP(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) error {
	return pb.RegisterNotificationServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
}

func (s *NotificationService) Send(ctx context.Context, req *pb.SendRequest) (*pb.SendResponse, error) {
	godotenv.Load(".env")
	pluginName := os.Getenv("NOTIFICATION_PLUGIN")

	if req.Plugin != "" {
		pluginName = req.Plugin
	}
	plugin, ok := s.registry.GetPlugin("notification", pluginName)
	if !ok {
		return nil, status.Error(codes.Unimplemented, "plugin not found")
	}
	return plugin.(NotificationServiceInterface).Send(ctx, req)
}

func (s *NotificationService) Delete(ctx context.Context, req *pb.DeleteRequest) (*pb.DeleteResponse, error) {
	godotenv.Load(".env")
	pluginName := os.Getenv("NOTIFICATION_PLUGIN")

	if req.Plugin != "" {
		pluginName = req.Plugin
	}
	plugin, ok := s.registry.GetPlugin("notification", pluginName)
	if !ok {
		return nil, status.Error(codes.Unimplemented, "plugin not found")
	}
	return plugin.(NotificationServiceInterface).Delete(ctx, req)
}

func (s *NotificationService) Update(ctx context.Context, req *pb.UpdateRequest) (*pb.UpdateResponse, error) {
	godotenv.Load(".env")
	pluginName := os.Getenv("NOTIFICATION_PLUGIN")

	if req.Plugin != "" {
		pluginName = req.Plugin
	}
	plugin, ok := s.registry.GetPlugin("notification", pluginName)
	if !ok {
		return nil, status.Error(codes.Unimplemented, "plugin not found")
	}
	return plugin.(NotificationServiceInterface).Update(ctx, req)
}
