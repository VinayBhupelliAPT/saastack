package interfaces

import (
	"context"
	"os"
	pb_payment "sample/proto/payment"

	"github.com/joho/godotenv"
)

const PaymentInterfaceName = "payment"

type PaymentPlugin interface {
	Charge(data map[string]string) string
	Refund(data map[string]string) string
	Status(data map[string]string) string
}

var PaymentPluginsRegistery = make(map[string]PaymentPlugin)

type PaymentServer struct {
	pb_payment.UnimplementedPaymentServiceServer
}

func NewPaymentServer() *PaymentServer {
	return &PaymentServer{}
}

func RegisterPaymentPlugins(pluginMap map[string]interface{}) {
	for name, plugin := range pluginMap {
		PaymentPluginsRegistery[name] = plugin.(PaymentPlugin)
	}
}

func (s *PaymentServer) Charge(ctx context.Context, req *pb_payment.ChargeRequest) (*pb_payment.ChargeResponse, error) {
	godotenv.Load(".env")
	defaultPlugin := os.Getenv("PAYMENT_PLUGIN")
	handler := PaymentPluginsRegistery[defaultPlugin]
	result := handler.Charge(map[string]string{"message": req.Message})
	return &pb_payment.ChargeResponse{Result: result}, nil
}

func (s *PaymentServer) Refund(ctx context.Context, req *pb_payment.RefundRequest) (*pb_payment.RefundResponse, error) {
	godotenv.Load(".env")
	defaultPlugin := os.Getenv("PAYMENT_PLUGIN")
	handler := PaymentPluginsRegistery[defaultPlugin]
	result := handler.Refund(map[string]string{"message": req.Message})
	return &pb_payment.RefundResponse{Result: result}, nil
}

func (s *PaymentServer) Status(ctx context.Context, req *pb_payment.StatusRequest) (*pb_payment.StatusResponse, error) {
	godotenv.Load(".env")
	defaultPlugin := os.Getenv("PAYMENT_PLUGIN")
	handler := PaymentPluginsRegistery[defaultPlugin]
	result := handler.Status(map[string]string{"message": req.Message})
	return &pb_payment.StatusResponse{Result: result}, nil
}
