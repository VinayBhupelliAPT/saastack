package core

import (
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

func GetGRPCServer() *grpc.Server {
	return grpc.NewServer()
}

// Start gRPC server
func StartGRPCServer(s *grpc.Server) {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	fmt.Println("gRPC server started on :50051")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

// Start HTTP gateway
func GetHTTPGateway() *runtime.ServeMux {
	return runtime.NewServeMux()
}

func StartHTTPGateway(mux *runtime.ServeMux) {
	fmt.Println("HTTP gateway started on :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("Failed to serve HTTP gateway: %v", err)
	}
}
