package plugins

import (
	"context"
	"fmt"
	"saastack/core"
	notification "saastack/interfaces/notification"
	pb_notification "saastack/interfaces/notification/proto"
)

// EmailPlugin implements the NotificationServiceInterface
type EmailPlugin struct {
	pb_notification.UnimplementedNotificationServiceServer
}

// NewEmailPlugin creates a new instance of EmailPlugin and verifies interface implementation
func NewEmailPlugin() notification.NotificationServiceInterface {
	plugin := &EmailPlugin{}
	// Verify interface implementation at compile time
	var _ notification.NotificationServiceInterface = plugin
	return plugin
}

func init() {
	// Register the plugin with the core registry
	plugin := NewEmailPlugin()
	core.GlobalRegistry.RegisterPlugin("notification", "email", plugin)
}

func (e *EmailPlugin) Send(ctx context.Context, req *pb_notification.SendRequest) (*pb_notification.SendResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("request cannot be nil")
	}
	msg := req.Message
	fmt.Println("EmailPlugin sending:", msg)
	return &pb_notification.SendResponse{Result: "EmailPlugin sent: " + msg}, nil
}

func (e *EmailPlugin) Delete(ctx context.Context, req *pb_notification.DeleteRequest) (*pb_notification.DeleteResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("request cannot be nil")
	}
	msg := req.Message
	fmt.Println("EmailPlugin deleting:", msg)
	return &pb_notification.DeleteResponse{Result: "EmailPlugin deleted: " + msg}, nil
}

func (e *EmailPlugin) Update(ctx context.Context, req *pb_notification.UpdateRequest) (*pb_notification.UpdateResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("request cannot be nil")
	}
	msg := req.Message
	fmt.Println("EmailPlugin updating:", msg)
	return &pb_notification.UpdateResponse{Result: "EmailPlugin updated: " + msg}, nil
}
