package payment

import (
	"context"
	"os"

	"saastack/core"
	pb "saastack/interfaces/payment/proto"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// PaymentServiceInterface defines the interface for payment service
type PaymentServiceInterface interface {
	pb.PaymentServiceServer
}

// Registry interface for plugin management
type PluginRegistry interface {
	GetPlugin(interfaceName, pluginName string) (interface{}, bool)
}

// PaymentService implements the PaymentServiceServer interface
type PaymentService struct {
	pb.UnimplementedPaymentServiceServer
	registry PluginRegistry
}

// NewPaymentService creates a new instance of PaymentService
func NewPaymentService(registry PluginRegistry) *PaymentService {
	return &PaymentService{
		registry: registry,
	}
}

func init() {
	service := NewPaymentService(core.GlobalRegistry)
	core.GlobalRegistry.RegisterService("payment", service)
}

// RegisterGRPC registers the service with gRPC server
func (s *PaymentService) RegisterGRPC(server *grpc.Server) {
	pb.RegisterPaymentServiceServer(server, s)
}

// RegisterHTTP registers the service with HTTP gateway
func (s *PaymentService) RegisterHTTP(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) error {
	return pb.RegisterPaymentServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
}

func (s *PaymentService) Charge(ctx context.Context, req *pb.ChargeRequest) (*pb.ChargeResponse, error) {
	godotenv.Load(".env")
	pluginName := os.Getenv("PAYMENT_PLUGIN")

	if req.Plugin != "" {
		pluginName = req.Plugin
	}
	plugin, ok := s.registry.GetPlugin("payment", pluginName)
	if !ok {
		return nil, status.Error(codes.Unimplemented, "plugin not found")
	}
	return plugin.(PaymentServiceInterface).Charge(ctx, req)
}

func (s *PaymentService) Refund(ctx context.Context, req *pb.RefundRequest) (*pb.RefundResponse, error) {
	godotenv.Load(".env")
	pluginName := os.Getenv("PAYMENT_PLUGIN")

	if req.Plugin != "" {
		pluginName = req.Plugin
	}
	plugin, ok := s.registry.GetPlugin("payment", pluginName)
	if !ok {
		return nil, status.Error(codes.Unimplemented, "plugin not found")
	}
	return plugin.(PaymentServiceInterface).Refund(ctx, req)
}

func (s *PaymentService) Status(ctx context.Context, req *pb.StatusRequest) (*pb.StatusResponse, error) {
	godotenv.Load(".env")
	pluginName := os.Getenv("PAYMENT_PLUGIN")

	if req.Plugin != "" {
		pluginName = req.Plugin
	}
	plugin, ok := s.registry.GetPlugin("payment", pluginName)
	if !ok {
		return nil, status.Error(codes.Unimplemented, "plugin not found")
	}
	return plugin.(PaymentServiceInterface).Status(ctx, req)
}
