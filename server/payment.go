package server

import (
	"context"
	"fmt"
	"sample/core"
	"sample/plugins"
	pb_payment "sample/proto/payment"
)

type PaymentServer struct {
	pb_payment.UnimplementedPaymentServiceServer
	pluginMap map[string]interface{}
}

func NewPaymentServer() *PaymentServer {
	return &PaymentServer{
		pluginMap: map[string]interface{}{
			"stripe": plugins.NewStripePlugin(),
		},
	}
}

func (s *PaymentServer) Charge(ctx context.Context, req *pb_payment.ChargeRequest) (*pb_payment.ChargeResponse, error) {
	_, ok := s.pluginMap[req.Plugin]
	if !ok {
		return nil, fmt.Errorf("plugin '%s' not found", req.Plugin)
	}

	handler := core.PluginRegistry["payment"][req.Plugin]["Charge"]
	if handler == nil {
		return nil, fmt.Errorf("method 'Charge' not found for plugin '%s'", req.Plugin)
	}

	result := handler(map[string]string{"message": req.Message})
	return &pb_payment.ChargeResponse{Result: result}, nil
}

func (s *PaymentServer) Refund(ctx context.Context, req *pb_payment.RefundRequest) (*pb_payment.RefundResponse, error) {
	_, ok := s.pluginMap[req.Plugin]
	if !ok {
		return nil, fmt.Errorf("plugin '%s' not found", req.Plugin)
	}

	handler := core.PluginRegistry["payment"][req.Plugin]["Refund"]
	if handler == nil {
		return nil, fmt.Errorf("method 'Refund' not found for plugin '%s'", req.Plugin)
	}

	result := handler(map[string]string{"message": req.Message})
	return &pb_payment.RefundResponse{Result: result}, nil
}

func (s *PaymentServer) Status(ctx context.Context, req *pb_payment.StatusRequest) (*pb_payment.StatusResponse, error) {
	_, ok := s.pluginMap[req.Plugin]
	if !ok {
		return nil, fmt.Errorf("plugin '%s' not found", req.Plugin)
	}

	handler := core.PluginRegistry["payment"][req.Plugin]["Status"]
	if handler == nil {
		return nil, fmt.Errorf("method 'Status' not found for plugin '%s'", req.Plugin)
	}

	result := handler(map[string]string{"message": req.Message})
	return &pb_payment.StatusResponse{Result: result}, nil
}
