package core

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func loggingUnaryInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	start := time.Now()

	if md, ok := metadata.FromIncomingContext(ctx); ok {
		log.Printf("Metadata: %v", md)
	}

	log.Printf("gRPC method: %s, Request: %+v", info.FullMethod, req)

	resp, err := handler(ctx, req)

	log.Printf("gRPC method: %s, Response: %+v, Duration: %s", info.FullMethod, resp, time.Since(start))
	if err != nil {
		log.Printf("gRPC error: %v", err)
	}

	return resp, err
}
func loggingHTTPMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		log.Printf("HTTP %s %s from %s", r.Method, r.URL.Path, r.RemoteAddr)
		next.ServeHTTP(w, r)
		log.Printf("HTTP Completed in %s", time.Since(start))
	})
}
func GetGRPCServer() *grpc.Server {
	return grpc.NewServer(grpc.UnaryInterceptor(loggingUnaryInterceptor))
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
	if err := http.ListenAndServe(":8080", loggingHTTPMiddleware(mux)); err != nil {
		log.Fatalf("Failed to serve HTTP gateway: %v", err)
	}
}
