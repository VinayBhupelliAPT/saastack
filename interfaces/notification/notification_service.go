package notification

import (
	"context"
	"fmt"
	"os"
	pb_notification "saastack/interfaces/notification/proto"

	"github.com/joho/godotenv"	
)

const InterfaceName = "notification"

var NotificationPluginsRegistery = make(map[string]NotificationPluginDetails)

type NotificationPlugin interface {
	Send(ctx context.Context, req *pb_notification.SendRequest) (*pb_notification.SendResponse, error)
	Delete(ctx context.Context, req *pb_notification.DeleteRequest) (*pb_notification.DeleteResponse, error)
	Update(ctx context.Context, req *pb_notification.UpdateRequest) (*pb_notification.UpdateResponse, error)
}
type NotificationService struct {
	pb_notification.UnimplementedNotificationServiceServer
}

func NewNotificationService() *NotificationService {
	return &NotificationService{}
}

type NotificationPluginDetails struct {
	Name   string
	Plugin NotificationPlugin
	Client pb_notification.NotificationServiceClient
}

func RegisterNotificationPlugin(details NotificationPluginDetails) {
	NotificationPluginsRegistery[details.Name] = details
}

func (s *NotificationService) Send(ctx context.Context, req *pb_notification.SendRequest) (*pb_notification.SendResponse, error) {
	godotenv.Load(".env")
	defaultPlugin := os.Getenv("NOTIFICATION_PLUGIN")
	var details NotificationPluginDetails
	if req.Plugin != "" {
		details = NotificationPluginsRegistery[req.Plugin]
	} else {
		details = NotificationPluginsRegistery[defaultPlugin]
	}

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

func (s *NotificationService) Delete(ctx context.Context, req *pb_notification.DeleteRequest) (*pb_notification.DeleteResponse, error) {
	godotenv.Load(".env")
	defaultPlugin := os.Getenv("NOTIFICATION_PLUGIN")
	var details NotificationPluginDetails
	if req.Plugin != "" {
		details = NotificationPluginsRegistery[req.Plugin]
	} else {
		details = NotificationPluginsRegistery[defaultPlugin]
	}

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

func (s *NotificationService) Update(ctx context.Context, req *pb_notification.UpdateRequest) (*pb_notification.UpdateResponse, error) {
	godotenv.Load(".env")
	defaultPlugin := os.Getenv("NOTIFICATION_PLUGIN")
	var details NotificationPluginDetails
	if req.Plugin != "" {
		details = NotificationPluginsRegistery[req.Plugin]
	} else {
		details = NotificationPluginsRegistery[defaultPlugin]
	}

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
