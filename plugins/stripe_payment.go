package plugins

import (
	"context"
	"fmt"
	pb_notification "sample/proto/notification"
	pb_payment "sample/proto/payment"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type StripePlugin struct {
	pb_payment.UnimplementedPaymentServiceServer
}

func NewStripePlugin() *StripePlugin {
	return &StripePlugin{}
}

func (s *StripePlugin) Charge(ctx context.Context, req *pb_payment.ChargeRequest) (*pb_payment.ChargeResponse, error) {
	msg := req.Message
	fmt.Println("StripePlugin charging:", msg)
	// using notification service
	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Printf("Failed to connect to notification service: %v\n", err)
	}

	client := pb_notification.NewNotificationServiceClient(conn)
	req1 := &pb_notification.SendRequest{
		Message: "Payment successful: " + msg,
	}
	client.Send(context.Background(), req1)
	return &pb_payment.ChargeResponse{Result: "StripePlugin charged: " + msg}, nil
}

func (s *StripePlugin) Refund(ctx context.Context, req *pb_payment.RefundRequest) (*pb_payment.RefundResponse, error) {
	msg := req.Message
	fmt.Println("StripePlugin refunding:", msg)

	// using notification service
	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Printf("Failed to connect to notification service: %v\n", err)
	}

	client := pb_notification.NewNotificationServiceClient(conn)
	req1 := &pb_notification.SendRequest{
		Message: "Refund processed: " + msg,
	}
	client.Send(context.Background(), req1)
	return &pb_payment.RefundResponse{Result: "StripePlugin refunded: " + msg}, nil
}

func (s *StripePlugin) Status(ctx context.Context, req *pb_payment.StatusRequest) (*pb_payment.StatusResponse, error) {
	msg := req.Message
	fmt.Println("StripePlugin checking status for:", msg)
	return &pb_payment.StatusResponse{Result: "StripePlugin status: Success for " + msg}, nil
}
