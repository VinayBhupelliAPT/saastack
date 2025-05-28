package interfaces

import (
	"context"
	"os"
	pb_notification "sample/proto/notification"

	"github.com/joho/godotenv"
)

const InterfaceName = "notification"

type NotificationPlugin interface {
	Send(data map[string]string) string
	Delete(data map[string]string) string
	Update(data map[string]string) string
}

var NotificationPluginsRegistery = make(map[string]NotificationPlugin)

type NotificationServer struct {
	pb_notification.UnimplementedNotificationServiceServer
}

func NewNotificationServer() *NotificationServer {
	return &NotificationServer{}
}

// instead of storing pluginMap store a single plugin
func RegisterNotificationPlugins(pluginMap map[string]interface{}) {
	for name, plugin := range pluginMap {
		NotificationPluginsRegistery[name] = plugin.(NotificationPlugin)
	}
}

func (s *NotificationServer) Send(ctx context.Context, req *pb_notification.SendRequest) (*pb_notification.SendResponse, error) {
	godotenv.Load(".env")
	defaultPlugin := os.Getenv("NOTIFICATION_PLUGIN")
	handler := NotificationPluginsRegistery[defaultPlugin]
	result := handler.Send(map[string]string{"message": req.Message})
	return &pb_notification.SendResponse{Result: result}, nil
}

func (s *NotificationServer) Delete(ctx context.Context, req *pb_notification.DeleteRequest) (*pb_notification.DeleteResponse, error) {
	godotenv.Load(".env")
	defaultPlugin := os.Getenv("NOTIFICATION_PLUGIN")
	handler := NotificationPluginsRegistery[defaultPlugin]
	result := handler.Delete(map[string]string{"message": req.Message})
	return &pb_notification.DeleteResponse{Result: result}, nil
}

func (s *NotificationServer) Update(ctx context.Context, req *pb_notification.UpdateRequest) (*pb_notification.UpdateResponse, error) {
	godotenv.Load(".env")
	defaultPlugin := os.Getenv("NOTIFICATION_PLUGIN")
	handler := NotificationPluginsRegistery[defaultPlugin]
	result := handler.Update(map[string]string{"message": req.Message})
	return &pb_notification.UpdateResponse{Result: result}, nil
}
