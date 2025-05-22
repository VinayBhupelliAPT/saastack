package interfaces

import "sample/core"

const InterfaceName = "notification"

type NotificationPlugin interface {
	Send(data map[string]string) string
	Delete(data map[string]string) string
	Update(data map[string]string) string
}

func init() {
	core.RegisterInterface(InterfaceName)
}

func RegisterNotificationPlugin(name string, plugin NotificationPlugin) {
	core.RegisterMethod(InterfaceName, name, "Send", plugin.Send)
	core.RegisterMethod(InterfaceName, name, "Delete", plugin.Delete)
	core.RegisterMethod(InterfaceName, name, "Update", plugin.Update)
}
