package interfaces

import (
	"context"
	"fmt"
	"os"
	pb_payment "sample/proto/payment"

	"github.com/joho/godotenv"
)

const PaymentInterfaceName = "payment"

type PaymentPlugin interface {
	Charge(ctx context.Context, req *pb_payment.ChargeRequest) (*pb_payment.ChargeResponse, error)
	Refund(ctx context.Context, req *pb_payment.RefundRequest) (*pb_payment.RefundResponse, error)
	Status(ctx context.Context, req *pb_payment.StatusRequest) (*pb_payment.StatusResponse, error)
}
type PaymentPluginDetails struct {
	Name   string
	Plugin PaymentPlugin
	Client pb_payment.PaymentServiceClient
}

var PaymentPluginsRegistery = make(map[string]PaymentPluginDetails)

type PaymentServer struct {
	pb_payment.UnimplementedPaymentServiceServer
}

func NewPaymentServer() *PaymentServer {
	return &PaymentServer{}
}

func RegisterPaymentPlugin(details PaymentPluginDetails) {
	PaymentPluginsRegistery[details.Name] = details
}

func (s *PaymentServer) Charge(ctx context.Context, req *pb_payment.ChargeRequest) (*pb_payment.ChargeResponse, error) {
	godotenv.Load(".env")
	defaultPlugin := os.Getenv("PAYMENT_PLUGIN")
	details := PaymentPluginsRegistery[defaultPlugin]
	if details.Plugin != nil {
		result, err := details.Plugin.Charge(ctx, req)
		if err != nil {
			return nil, err
		}
		return result, nil
	}

	if details.Client != nil {
		return details.Client.Charge(ctx, req)
	}

	return nil, fmt.Errorf("no valid plugin or client found for plugin: %s", defaultPlugin)
}

func (s *PaymentServer) Refund(ctx context.Context, req *pb_payment.RefundRequest) (*pb_payment.RefundResponse, error) {
	godotenv.Load(".env")
	defaultPlugin := os.Getenv("PAYMENT_PLUGIN")
	details := PaymentPluginsRegistery[defaultPlugin]
	if details.Plugin != nil {
		result, err := details.Plugin.Refund(ctx, req)
		if err != nil {
			return nil, err
		}
		return result, nil
	}

	if details.Client != nil {
		return details.Client.Refund(ctx, req)
	}

	return nil, fmt.Errorf("no valid plugin or client found for plugin: %s", defaultPlugin)
}

func (s *PaymentServer) Status(ctx context.Context, req *pb_payment.StatusRequest) (*pb_payment.StatusResponse, error) {
	godotenv.Load(".env")
	defaultPlugin := os.Getenv("PAYMENT_PLUGIN")
	details := PaymentPluginsRegistery[defaultPlugin]
	if details.Plugin != nil {
		result, err := details.Plugin.Status(ctx, req)
		if err != nil {
			return nil, err
		}
		return result, nil
	}

	if details.Client != nil {
		return details.Client.Status(ctx, req)
	}

	return nil, fmt.Errorf("no valid plugin or client found for plugin: %s", defaultPlugin)
}
