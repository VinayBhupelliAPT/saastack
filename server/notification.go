package server

import (
	"context"
	"fmt"
	"sample/core"
	"sample/plugins"
	pb_notification "sample/proto/notification"
)

type NotificationServer struct {
	pb_notification.UnimplementedNotificationServiceServer
	pluginMap map[string]interface{}
}

func NewNotificationServer() *NotificationServer {
	return &NotificationServer{
		pluginMap: map[string]interface{}{
			"email": &plugins.EmailPlugin{},
		},
	}
}

func (s *NotificationServer) Send(ctx context.Context, req *pb_notification.SendRequest) (*pb_notification.SendResponse, error) {
	_, ok := s.pluginMap[req.Plugin]
	if !ok {
		return nil, fmt.Errorf("plugin '%s' not found", req.Plugin)
	}

	handler := core.PluginRegistry["notification"][req.Plugin]["Send"]
	if handler == nil {
		return nil, fmt.Errorf("method 'Send' not found for plugin '%s'", req.Plugin)
	}

	result := handler(map[string]string{"message": req.Message})
	return &pb_notification.SendResponse{Result: result}, nil
}

func (s *NotificationServer) Delete(ctx context.Context, req *pb_notification.DeleteRequest) (*pb_notification.DeleteResponse, error) {
	_, ok := s.pluginMap[req.Plugin]
	if !ok {
		return nil, fmt.Errorf("plugin '%s' not found", req.Plugin)
	}

	handler := core.PluginRegistry["notification"][req.Plugin]["Delete"]
	if handler == nil {
		return nil, fmt.Errorf("method 'Delete' not found for plugin '%s'", req.Plugin)
	}

	result := handler(map[string]string{"message": req.Message})
	return &pb_notification.DeleteResponse{Result: result}, nil
}

func (s *NotificationServer) Update(ctx context.Context, req *pb_notification.UpdateRequest) (*pb_notification.UpdateResponse, error) {
	_, ok := s.pluginMap[req.Plugin]
	if !ok {
		return nil, fmt.Errorf("plugin '%s' not found", req.Plugin)
	}

	handler := core.PluginRegistry["notification"][req.Plugin]["Update"]
	if handler == nil {
		return nil, fmt.Errorf("method 'Update' not found for plugin '%s'", req.Plugin)
	}

	result := handler(map[string]string{"message": req.Message})
	return &pb_notification.UpdateResponse{Result: result}, nil
}
