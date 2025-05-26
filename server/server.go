package server

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"sample/core"
	pb_notification "sample/proto/notification"
	pb_payment "sample/proto/payment"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

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
