package interfaces

import (
	"context"
	"fmt"
	"os"
	pb_notification "sample/proto/notification"

	"github.com/joho/godotenv"
)

const InterfaceName = "notification"

var NotificationPluginsRegistery = make(map[string]NotificationPluginDetails)
type NotificationPlugin interface {
	Send(ctx context.Context, req *pb_notification.SendRequest) (*pb_notification.SendResponse, error)
	Delete(ctx context.Context, req *pb_notification.DeleteRequest) (*pb_notification.DeleteResponse, error)
	Update(ctx context.Context, req *pb_notification.UpdateRequest) (*pb_notification.UpdateResponse, error)
}
type NotificationServer struct {
	pb_notification.UnimplementedNotificationServiceServer
}

func NewNotificationServer() *NotificationServer {
	return &NotificationServer{}
}

type NotificationPluginDetails struct {
	Name   string
	Plugin NotificationPlugin
	Client pb_notification.NotificationServiceClient
}

// instead of storing pluginMap store a single plugin
func RegisterNotificationPlugin(details NotificationPluginDetails) {
	NotificationPluginsRegistery[details.Name] = details
}

func (s *NotificationServer) Send(ctx context.Context, req *pb_notification.SendRequest) (*pb_notification.SendResponse, error) {
	godotenv.Load(".env")
	defaultPlugin := os.Getenv("NOTIFICATION_PLUGIN")
	details := NotificationPluginsRegistery[defaultPlugin]

	if details.Plugin != nil {
		result, err := details.Plugin.Send(ctx, req)
		if err != nil {
			return nil, err
		}
		fmt.Println(result)
		return result, nil
	}

	if details.Client != nil {
		return details.Client.Send(ctx, req)
	}

	return nil, fmt.Errorf("no valid plugin or client found for plugin: %s", defaultPlugin)
}

func (s *NotificationServer) Delete(ctx context.Context, req *pb_notification.DeleteRequest) (*pb_notification.DeleteResponse, error) {
	godotenv.Load(".env")
	defaultPlugin := os.Getenv("NOTIFICATION_PLUGIN")
	details := NotificationPluginsRegistery[defaultPlugin]

	if details.Plugin != nil {
		result, err := details.Plugin.Delete(ctx, req)
		if err != nil {
			return nil, err
		}
		return result, nil
	}

	if details.Client != nil {
		return details.Client.Delete(ctx, req)
	}

	return nil, fmt.Errorf("no valid plugin or client found for plugin: %s", defaultPlugin)
}

func (s *NotificationServer) Update(ctx context.Context, req *pb_notification.UpdateRequest) (*pb_notification.UpdateResponse, error) {
	godotenv.Load(".env")
	defaultPlugin := os.Getenv("NOTIFICATION_PLUGIN")
	details := NotificationPluginsRegistery[defaultPlugin]
	if details.Plugin != nil {
		result, err := details.Plugin.Update(ctx, req)
		if err != nil {
			return nil, err
		}
		return result, nil
	}

	if details.Client != nil {
		return details.Client.Update(ctx, req)
	}

	return nil, fmt.Errorf("no valid plugin or client found for plugin: %s", defaultPlugin)

}
