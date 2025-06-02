package plugins

import (
	"context"
	"fmt"
	pb_payment "saastack/interfaces/payment/proto"
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
	return &pb_payment.ChargeResponse{Result: "StripePlugin charged: " + msg}, nil
}

func (s *StripePlugin) Refund(ctx context.Context, req *pb_payment.RefundRequest) (*pb_payment.RefundResponse, error) {
	msg := req.Message
	fmt.Println("StripePlugin refunding:", msg)
	return &pb_payment.RefundResponse{Result: "StripePlugin refunded: " + msg}, nil
}

func (s *StripePlugin) Status(ctx context.Context, req *pb_payment.StatusRequest) (*pb_payment.StatusResponse, error) {
	msg := req.Message
	fmt.Println("StripePlugin checking status for:", msg)
	return &pb_payment.StatusResponse{Result: "StripePlugin status: Success for " + msg}, nil
}
