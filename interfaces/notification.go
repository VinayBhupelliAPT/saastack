package interfaces

const InterfaceName = "notification"

type NotificationPlugin interface {
	Send(data map[string]string) string
	Delete(data map[string]string) string
	Update(data map[string]string) string
}
