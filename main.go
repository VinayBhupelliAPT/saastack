package main

import (
	"context"
	"fmt"
	"os"
	"sample/core"
	"sample/interfaces"
	"sample/plugins"
	pb_notification "sample/proto/notification"
	pb_payment "sample/proto/payment"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"gopkg.in/yaml.v3"
)

// Register Interfaces
var InterfaceRegistry = make(map[string]bool)

func RegisterInterfaces(interfaceList []string) {
	for _, iface := range interfaceList {
		InterfaceRegistry[iface] = true
	}
}

type Interface struct {
	Name string `yaml:"name"`
}

type InterfaceConfig struct {
	Interfaces []Interface `yaml:"interfaces"`
}

type Plugin struct {
	Name      string `yaml:"name"`
	Interface string `yaml:"interface"`
	Instance  string `yaml:"instance"`
}

type PluginConfig struct {
	Plugins []Plugin `yaml:"plugins"`
}

func LoadInterfaceConfig(filename string) (*InterfaceConfig, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var config InterfaceConfig
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
func LoadPluginConfig(filename string) (*PluginConfig, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var config PluginConfig
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

func main() {
	// Register Interfaces in core
	interfaceConfig, err := LoadInterfaceConfig("config/interfaces.yaml")
	if err != nil {
		fmt.Printf("Error loading interface config: %v\n", err)
		return
	}
	interfaceList := make([]string, 0)
	for _, iface := range interfaceConfig.Interfaces {
		interfaceList = append(interfaceList, iface.Name)
	}
	RegisterInterfaces(interfaceList)

	// Register Plugins in interfaces
	pluginConfig, err := LoadPluginConfig("config/plugins.yaml")
	if err != nil {
		fmt.Printf("Error loading plugin config: %v\n", err)
		return
	}
	notificationPluginsMap := make(map[string]interface{})
	paymentPluginsMap := make(map[string]interface{})

	for _, p := range pluginConfig.Plugins {
		if !InterfaceRegistry[p.Interface] {
			fmt.Printf("Interface %s not found in interface registry\n", p.Interface)
			continue
		}
		switch p.Instance {
		case "NewEmailPlugin":
			notificationPluginsMap[p.Name] = plugins.NewEmailPlugin()
		case "NewStripePlugin":
			paymentPluginsMap[p.Name] = plugins.NewStripePlugin()
		default:
			fmt.Printf("Unknown plugin instance: %s", p.Instance)
		}
	}
	interfaces.RegisterNotificationPlugins(notificationPluginsMap)
	interfaces.RegisterPaymentPlugins(paymentPluginsMap)

	// Register Services
	notificationServer := interfaces.NewNotificationServer()
	paymentServer := interfaces.NewPaymentServer()
	s := core.GetGRPCServer()
	pb_notification.RegisterNotificationServiceServer(s, notificationServer)
	pb_payment.RegisterPaymentServiceServer(s, paymentServer)

	// Start the gRPC server from core
	go core.StartGRPCServer(s)

	// Start the HTTP gateway from core
	mux := core.GetHTTPGateway()
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	pb_notification.RegisterNotificationServiceHandlerFromEndpoint(
		ctx, mux, "localhost:50051", opts,
	)
	pb_payment.RegisterPaymentServiceHandlerFromEndpoint(
		ctx, mux, "localhost:50051", opts,
	)
	core.StartHTTPGateway(mux)

}
