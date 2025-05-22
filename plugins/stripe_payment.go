package plugins

import (
	"fmt"
	"sample/interfaces"
)

type StripePlugin struct{}

func (s *StripePlugin) Charge(data map[string]string) string {
	msg := data["message"]
	fmt.Println("StripePlugin charging:", msg)
	return "StripePlugin charged: " + msg
}

func (s *StripePlugin) Refund(data map[string]string) string {
	msg := data["message"]
	fmt.Println("StripePlugin refunding:", msg)
	return "StripePlugin refunded: " + msg
}

func (s *StripePlugin) Status(data map[string]string) string {
	msg := data["message"]
	fmt.Println("StripePlugin checking status for:", msg)
	return "StripePlugin status: Success for " + msg
}

func init() {
	plugin := &StripePlugin{}
	interfaces.RegisterPaymentPlugin("stripe", plugin)
}
