package server

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"sample/core"
	"sample/plugins"
	pb_notification "sample/proto/notification"
	pb_payment "sample/proto/payment"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type NotificationServer struct {
	pb_notification.UnimplementedNotificationServiceServer
	pluginMap map[string]interface{}
}

type PaymentServer struct {
	pb_payment.UnimplementedPaymentServiceServer
	pluginMap map[string]interface{}
}

func NewNotificationServer() *NotificationServer {
	return &NotificationServer{
		pluginMap: map[string]interface{}{
			"email": &plugins.EmailPlugin{},
		},
	}
}

func NewPaymentServer() *PaymentServer {
	return &PaymentServer{
		pluginMap: map[string]interface{}{
			"stripe": plugins.NewStripePlugin(),
		},
	}
}

func (s *NotificationServer) Send(ctx context.Context, req *pb_notification.SendRequest) (*pb_notification.SendResponse, error) {
	_, ok := s.pluginMap[req.Plugin]
	if !ok {
		return nil, fmt.Errorf("plugin '%s' not found", req.Plugin)
	}

	handler := core.PluginRegistry["notification"][req.Plugin]["Send"]
	if handler == nil {
		return nil, fmt.Errorf("method 'Send' not found for plugin '%s'", req.Plugin)
	}

	result := handler(map[string]string{"message": req.Message})
	return &pb_notification.SendResponse{Result: result}, nil
}

func (s *NotificationServer) Delete(ctx context.Context, req *pb_notification.DeleteRequest) (*pb_notification.DeleteResponse, error) {
	_, ok := s.pluginMap[req.Plugin]
	if !ok {
		return nil, fmt.Errorf("plugin '%s' not found", req.Plugin)
	}

	handler := core.PluginRegistry["notification"][req.Plugin]["Delete"]
	if handler == nil {
		return nil, fmt.Errorf("method 'Delete' not found for plugin '%s'", req.Plugin)
	}

	result := handler(map[string]string{"message": req.Message})
	return &pb_notification.DeleteResponse{Result: result}, nil
}

func (s *NotificationServer) Update(ctx context.Context, req *pb_notification.UpdateRequest) (*pb_notification.UpdateResponse, error) {
	_, ok := s.pluginMap[req.Plugin]
	if !ok {
		return nil, fmt.Errorf("plugin '%s' not found", req.Plugin)
	}

	handler := core.PluginRegistry["notification"][req.Plugin]["Update"]
	if handler == nil {
		return nil, fmt.Errorf("method 'Update' not found for plugin '%s'", req.Plugin)
	}

	result := handler(map[string]string{"message": req.Message})
	return &pb_notification.UpdateResponse{Result: result}, nil
}

func (s *PaymentServer) Charge(ctx context.Context, req *pb_payment.ChargeRequest) (*pb_payment.ChargeResponse, error) {
	_, ok := s.pluginMap[req.Plugin]
	if !ok {
		return nil, fmt.Errorf("plugin '%s' not found", req.Plugin)
	}

	handler := core.PluginRegistry["payment"][req.Plugin]["Charge"]
	if handler == nil {
		return nil, fmt.Errorf("method 'Charge' not found for plugin '%s'", req.Plugin)
	}

	result := handler(map[string]string{"message": req.Message})
	return &pb_payment.ChargeResponse{Result: result}, nil
}

func (s *PaymentServer) Refund(ctx context.Context, req *pb_payment.RefundRequest) (*pb_payment.RefundResponse, error) {
	_, ok := s.pluginMap[req.Plugin]
	if !ok {
		return nil, fmt.Errorf("plugin '%s' not found", req.Plugin)
	}

	handler := core.PluginRegistry["payment"][req.Plugin]["Refund"]
	if handler == nil {
		return nil, fmt.Errorf("method 'Refund' not found for plugin '%s'", req.Plugin)
	}

	result := handler(map[string]string{"message": req.Message})
	return &pb_payment.RefundResponse{Result: result}, nil
}

func (s *PaymentServer) Status(ctx context.Context, req *pb_payment.StatusRequest) (*pb_payment.StatusResponse, error) {
	_, ok := s.pluginMap[req.Plugin]
	if !ok {
		return nil, fmt.Errorf("plugin '%s' not found", req.Plugin)
	}

	handler := core.PluginRegistry["payment"][req.Plugin]["Status"]
	if handler == nil {
		return nil, fmt.Errorf("method 'Status' not found for plugin '%s'", req.Plugin)
	}

	result := handler(map[string]string{"message": req.Message})
	return &pb_payment.StatusResponse{Result: result}, nil
}

func StartServers() error {
	// Create plugin maps
	notificationServer := NewNotificationServer()
	paymentServer := NewPaymentServer()

	// Combine plugin maps
	pluginMap := make(map[string]interface{})
	for name, plugin := range notificationServer.pluginMap {
		pluginMap[name] = plugin
	}
	for name, plugin := range paymentServer.pluginMap {
		pluginMap[name] = plugin
	}

	// Initialize plugins from config
	if err := core.InitializeFromConfig("config/plugins.yaml", pluginMap); err != nil {
		return fmt.Errorf("failed to initialize plugins from config: %v", err)
	}

	// Start gRPC servers
	// Start notification gRPC server
	go func() {
		lis, err := net.Listen("tcp", ":50051")
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}
		s := grpc.NewServer()
		pb_notification.RegisterNotificationServiceServer(s, notificationServer)
		fmt.Println("Notification gRPC server started on :50051")
		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	// Start payment gRPC server
	go func() {
		lis, err := net.Listen("tcp", ":50052")
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}
		s := grpc.NewServer()
		pb_payment.RegisterPaymentServiceServer(s, paymentServer)
		fmt.Println("Payment gRPC server started on :50052")
		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	// Start HTTP gateway
	go func() {
		ctx := context.Background()
		ctx, cancel := context.WithCancel(ctx)
		defer cancel()

		mux := runtime.NewServeMux()
		opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

		// Register notification service
		if err := pb_notification.RegisterNotificationServiceHandlerFromEndpoint(
			ctx, mux, "localhost:50051", opts,
		); err != nil {
			log.Fatalf("Failed to register notification gateway: %v", err)
		}

		// Register payment service
		if err := pb_payment.RegisterPaymentServiceHandlerFromEndpoint(
			ctx, mux, "localhost:50052", opts,
		); err != nil {
			log.Fatalf("Failed to register payment gateway: %v", err)
		}

		fmt.Println("HTTP gateway started on :8080")
		if err := http.ListenAndServe(":8080", mux); err != nil {
			log.Fatalf("Failed to serve HTTP gateway: %v", err)
		}
	}()

	// Block forever
	select {}
}
