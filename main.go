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

type Interface struct {
	Name string `yaml:"name"`
}

type InterfaceConfig struct {
	Interfaces []Interface `yaml:"interfaces"`
}

type Plugin struct {
	Name       string `yaml:"name"`
	Interface  string `yaml:"interface"`
	Instance   string `yaml:"instance"`
	Deployment string `yaml:"deployment"`
	Source     string `yaml:"source"`
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

// Register Interfaces
var InterfaceRegistry = make(map[string]bool)

func main() {
	interfaceConfig, err := LoadInterfaceConfig("config/interfaces.yaml")
	if err != nil {
		fmt.Printf("Error loading interface config: %v\n", err)
		return
	}
	pluginConfig, err := LoadPluginConfig("config/plugins.yaml")
	if err != nil {
		fmt.Printf("Error loading plugin config: %v\n", err)
		return
	}
	// get gRPC and HTTP gateway
	s := core.GetGRPCServer()
	mux := core.GetHTTPGateway()
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	// Register Interfaces
	interfaceList := make(map[string]bool)
	for _, iface := range interfaceConfig.Interfaces {
		interfaceList[iface.Name] = true
		// Register Services
		switch iface.Name {
		case "notification":
			notificationServer := interfaces.NewNotificationServer()
			pb_notification.RegisterNotificationServiceServer(s, notificationServer)
			pb_notification.RegisterNotificationServiceHandlerFromEndpoint(
				ctx, mux, "localhost:50051", opts,
			)
		case "payment":
			paymentServer := interfaces.NewPaymentServer()
			pb_payment.RegisterPaymentServiceServer(s, paymentServer)
			pb_payment.RegisterPaymentServiceHandlerFromEndpoint(
				ctx, mux, "localhost:50051", opts,
			)
		}
	}

	// Register Plugins in interfaces
	for _, p := range pluginConfig.Plugins {
		switch p.Interface {
		case "notification":
			registerNotification(p)
		case "payment":
			registerPayment(p)
		default:
			fmt.Printf("Unknown plugin instance: %s", p.Instance)
		}
	}

	// Start the HTTP gateway from core
	go core.StartHTTPGateway(mux)
	// Start the gRPC server
	core.StartGRPCServer(s)
}
func registerNotification(p Plugin) {
	var data interfaces.NotificationPluginDetails
	data.Name = p.Name
	if p.Deployment == "Microservice" {
		if p.Name == "email" {
			go func() {
				plugins.StartEmailNotificationMicroservice("50052")
			}()
		}
		conn, err := grpc.NewClient(p.Source, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			fmt.Printf("Failed to connect to notification service: %v\n", err)
			return
		}
		client := pb_notification.NewNotificationServiceClient(conn)
		data.Client = client
	} else {
		data.Plugin = plugins.NewEmailPlugin()
	}
	interfaces.RegisterNotificationPlugin(data)
}
func registerPayment(p Plugin) {
	var data interfaces.PaymentPluginDetails
	data.Name = p.Name
	if p.Deployment == "Microservice" {
		conn, err := grpc.NewClient(p.Source, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			fmt.Printf("Failed to connect to payment service: %v\n", err)
			return
		}
		client := pb_payment.NewPaymentServiceClient(conn)
		data.Client = client
	} else {
		data.Plugin = plugins.NewStripePlugin()
	}
	interfaces.RegisterPaymentPlugin(data)
}
