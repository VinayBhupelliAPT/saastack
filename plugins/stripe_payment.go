package plugins

import (
	"context"
	"fmt"
	"log"

	"saastack/core"
	notification_pb "saastack/interfaces/notification/proto"
	pb "saastack/interfaces/payment/proto"
)

// StripePlugin implements the PaymentServiceServer interface
type StripePlugin struct {
	pb.UnimplementedPaymentServiceServer
}

func NewStripePlugin() *StripePlugin {
	return &StripePlugin{}
}

func (s *StripePlugin) Charge(ctx context.Context, req *pb.ChargeRequest) (*pb.ChargeResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("request cannot be nil")
	}
	msg := req.Message
	fmt.Println("StripePlugin charging:", msg)

	// Send notification about the charge
	if notificationService, exists := core.GlobalRegistry.GetService("notification"); exists {
		notificationReq := &notification_pb.SendRequest{
			Message: fmt.Sprintf("Payment charged: %s", msg),
			Plugin:  "email",
		}
		if _, err := notificationService.(notification_pb.NotificationServiceServer).Send(ctx, notificationReq); err != nil {
			log.Printf("Failed to send charge notification: %v", err)
		}
	}

	return &pb.ChargeResponse{Result: "StripePlugin charged: " + msg}, nil
}

func (s *StripePlugin) Refund(ctx context.Context, req *pb.RefundRequest) (*pb.RefundResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("request cannot be nil")
	}
	msg := req.Message
	fmt.Println("StripePlugin refunding:", msg)

	// Send notification about the refund
	if notificationService, exists := core.GlobalRegistry.GetService("notification"); exists {
		notificationReq := &notification_pb.SendRequest{
			Message: fmt.Sprintf("Payment refunded: %s", msg),
			Plugin:  "email",
		}
		if _, err := notificationService.(notification_pb.NotificationServiceServer).Send(ctx, notificationReq); err != nil {
			log.Printf("Failed to send refund notification: %v", err)
		}
	}

	return &pb.RefundResponse{Result: "StripePlugin refunded: " + msg}, nil
}

func (s *StripePlugin) Status(ctx context.Context, req *pb.StatusRequest) (*pb.StatusResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("request cannot be nil")
	}
	msg := req.Message
	fmt.Println("StripePlugin checking status for:", msg)
	return &pb.StatusResponse{Result: "StripePlugin status: Success for " + msg}, nil
}

func init() {
	plugin := NewStripePlugin()
	core.GlobalRegistry.RegisterPlugin("payment", "stripe", plugin)
}
