package core

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"gopkg.in/yaml.v3"
)

// Registry part:
type ServiceRegistrar interface {
	RegisterGRPC(server *grpc.Server)
	RegisterHTTP(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) error
}

type Registry struct {
	plugins  map[string]map[string]interface{} // interface -> plugin_name -> plugin_instance
	services map[string]ServiceRegistrar       // service_name -> service_instance
}

func NewRegistry() *Registry {
	return &Registry{
		plugins:  make(map[string]map[string]interface{}),
		services: make(map[string]ServiceRegistrar),
	}
}

func (r *Registry) RegisterPlugin(interfaceName, pluginName string, plugin interface{}) {
	if _, exists := r.plugins[interfaceName]; !exists {
		r.plugins[interfaceName] = make(map[string]interface{})
	}
	r.plugins[interfaceName][pluginName] = plugin
}

func (r *Registry) GetPlugin(interfaceName, pluginName string) (interface{}, bool) {
	if plugins, exists := r.plugins[interfaceName]; exists {
		plugin, ok := plugins[pluginName]
		return plugin, ok
	}
	return nil, false
}

func (r *Registry) RegisterService(name string, service ServiceRegistrar) {
	r.services[name] = service
}

func (r *Registry) GetService(name string) (ServiceRegistrar, bool) {
	service, exists := r.services[name]
	return service, exists
}

func (r *Registry) GetAllServices() map[string]ServiceRegistrar {
	return r.services
}

var GlobalRegistry = NewRegistry()

// Config Part:
type Interface struct {
	Name string `yaml:"name"`
}

type Plugin struct {
	Name       string `yaml:"name"`
	Interface  string `yaml:"interface"`
	Instance   string `yaml:"instance"`
	Deployment string `yaml:"deployment"`
	Source     string `yaml:"source,omitempty"`
}

type Config struct {
	Interfaces []Interface `yaml:"interfaces"`
	Plugins    []Plugin    `yaml:"plugins"`
}

func LoadConfig(filename string) (*Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

var Configuration *Config

func InitializefromConfig(filename string) {
	config, err := LoadConfig(filename)
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		return
	}
	Configuration = config
}

// Logging Part:
func loggingUnaryInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	start := time.Now()

	// Log metadata
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		log.Println("Metadata:")
		for key, values := range md {
			log.Printf("  %s: %v", key, values)
		}
	}

	// Pretty-print request
	log.Printf("gRPC method: %s", info.FullMethod)
	if reqJSON, err := json.MarshalIndent(req, "", "  "); err == nil {
		log.Println("Request:\n", string(reqJSON))
	} else {
		log.Printf("Request: %+v", req)
	}

	// Call the handler
	resp, err := handler(ctx, req)

	// Pretty-print response
	if respJSON, err := json.MarshalIndent(resp, "", "  "); err == nil {
		log.Println("Response:\n", string(respJSON))
	} else {
		log.Printf("Response: %+v", resp)
	}

	log.Printf("Duration: %s", time.Since(start))
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

func GetHTTPGateway() *runtime.ServeMux {
	return runtime.NewServeMux()
}

func Start() {
	s := GetGRPCServer()
	mux := GetHTTPGateway()
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	log.Printf("GlobalRegistry has %d services registered", len(GlobalRegistry.services))

	// Register all services
	servicesRegistered := 0
	for name, service := range GlobalRegistry.GetAllServices() {
		log.Printf("Registering service: %s", name)

		// Register gRPC service
		service.RegisterGRPC(s)
		log.Printf("gRPC registration completed for %s", name)

		// Register HTTP handler
		opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
		if err := service.RegisterHTTP(ctx, mux, "localhost:50051", opts); err != nil {
			log.Printf("Failed to register HTTP handler for %s: %v", name, err)
		} else {
			log.Printf("HTTP registration completed for %s", name)
		}

		servicesRegistered++
	}

	log.Printf("Total services registered: %d", servicesRegistered)
	info := s.GetServiceInfo()
	pretty, err := json.MarshalIndent(info, "", "  ")
	if err != nil {
		log.Fatalf("Failed to marshal service info: %v", err)
	}
	log.Println("gRPC server info:\n", string(pretty))
	// Print information about registered gRPC services
	// log.Printf("gRPC server info: %+v", s.GetServiceInfo())

	// Start gRPC server
	go func() {
		lis, err := net.Listen("tcp", ":50051")
		if err != nil {
			log.Printf("failed to listen on gRPC port: %v", err)
			return
		}
		log.Println("gRPC server started on :50051")
		if err := s.Serve(lis); err != nil {
			log.Printf("failed to serve gRPC: %v", err)
		}
	}()

	// Start HTTP gateway
	log.Println("Starting HTTP gateway on :8080")
	if err := http.ListenAndServe(":8080", loggingHTTPMiddleware(mux)); err != nil {
		fmt.Printf("failed to start HTTP gateway on any port: %v", err)
	}
}
