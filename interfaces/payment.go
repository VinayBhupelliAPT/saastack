package interfaces

const PaymentInterfaceName = "payment"

type PaymentPlugin interface {
	Charge(data map[string]string) string
	Refund(data map[string]string) string
	Status(data map[string]string) string
}
