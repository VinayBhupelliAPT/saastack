package plugins

import (
	"context"
	"fmt"
	pb_notification "saastack/interfaces/notification/proto"
)

type EmailPlugin struct {
	pb_notification.UnimplementedNotificationServiceServer
}

func NewEmailPlugin() *EmailPlugin {
	return &EmailPlugin{}
}

func (e *EmailPlugin) Send(ctx context.Context, req *pb_notification.SendRequest) (*pb_notification.SendResponse, error) {
	msg := req.Message
	fmt.Println("EmailPlugin sending:", msg)
	return &pb_notification.SendResponse{Result: "EmailPlugin sent: " + msg}, nil
}

func (e *EmailPlugin) Delete(ctx context.Context, req *pb_notification.DeleteRequest) (*pb_notification.DeleteResponse, error) {
	msg := req.Message
	fmt.Println("EmailPlugin deleting:", msg)
	return &pb_notification.DeleteResponse{Result: "EmailPlugin deleted: " + msg}, nil
}

func (e *EmailPlugin) Update(ctx context.Context, req *pb_notification.UpdateRequest) (*pb_notification.UpdateResponse, error) {
	msg := req.Message
	fmt.Println("EmailPlugin updating:", msg)
	return &pb_notification.UpdateResponse{Result: "EmailPlugin updated: " + msg}, nil
}
