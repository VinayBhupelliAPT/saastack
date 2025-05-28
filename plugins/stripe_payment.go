package plugins

import (
	"context"
	"fmt"
	pb_notification "sample/proto/notification"

	"google.golang.org/grpc"
)

type StripePlugin struct {
	notificationClient pb_notification.NotificationServiceClient
}

func NewStripePlugin() *StripePlugin {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		fmt.Printf("Failed to connect to notification service: %v\n", err)
		return &StripePlugin{}
	}

	client := pb_notification.NewNotificationServiceClient(conn)
	return &StripePlugin{
		notificationClient: client,
	}
}

func (s *StripePlugin) Charge(data map[string]string) string {
	msg := data["message"]
	fmt.Println("StripePlugin charging:", msg)

	if s.notificationClient != nil {
		req := &pb_notification.SendRequest{
			Message: "Payment successful: " + msg,
		}
		_, err := s.notificationClient.Send(context.Background(), req)
		if err != nil {
			fmt.Printf("Failed to send notification: %v\n", err)
		}
	}

	return "StripePlugin charged: " + msg
}

func (s *StripePlugin) Refund(data map[string]string) string {
	msg := data["message"]
	fmt.Println("StripePlugin refunding:", msg)

	if s.notificationClient != nil {
		req := &pb_notification.SendRequest{
			Message: "Refund processed: " + msg,
		}
		_, err := s.notificationClient.Send(context.Background(), req)
		if err != nil {
			fmt.Printf("Failed to send notification: %v\n", err)
		}
	}

	return "StripePlugin refunded: " + msg
}

func (s *StripePlugin) Status(data map[string]string) string {
	msg := data["message"]
	fmt.Println("StripePlugin checking status for:", msg)
	return "StripePlugin status: Success for " + msg
}
