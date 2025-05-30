package plugins

import (
	"context"
	"fmt"
	"log"
	"net"
	pb_notification "sample/proto/notification"
	
	"google.golang.org/grpc"
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

func StartEmailNotificationMicroservice(port string) {
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("Failed to listen on port %s: %v", port, err)
	}

	s := grpc.NewServer()
	emailServer := NewEmailPlugin()
	pb_notification.RegisterNotificationServiceServer(s, emailServer)

	fmt.Printf("Email Notification microservice started on port %s\n", port)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve email notification microservice: %v", err)
	}
}
