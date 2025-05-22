package interfaces

import "sample/core"

const PaymentInterfaceName = "payment"

type PaymentPlugin interface {
	Charge(data map[string]string) string
	Refund(data map[string]string) string
	Status(data map[string]string) string
}

func init() {
	core.RegisterInterface(PaymentInterfaceName)
}

func RegisterPaymentPlugin(name string, plugin PaymentPlugin) {
	core.RegisterMethod(PaymentInterfaceName, name, "Charge", plugin.Charge)
	core.RegisterMethod(PaymentInterfaceName, name, "Refund", plugin.Refund)
	core.RegisterMethod(PaymentInterfaceName, name, "Status", plugin.Status)
}
