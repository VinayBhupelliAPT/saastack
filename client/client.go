package client

import (
	"context"
	"fmt"
	"time"

	pb_notification "sample/proto/notification"
	pb_payment "sample/proto/payment"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	notificationClient pb_notification.NotificationServiceClient
	paymentClient      pb_payment.PaymentServiceClient
}

func NewClient() (*Client, error) {
	// Connect to notification server
	notificationConn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to notification server: %v", err)
	}

	// Connect to payment server
	paymentConn, err := grpc.Dial("localhost:50052", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to payment server: %v", err)
	}

	return &Client{
		notificationClient: pb_notification.NewNotificationServiceClient(notificationConn),
		paymentClient:      pb_payment.NewPaymentServiceClient(paymentConn),
	}, nil
}

// Notification methods
func (c *Client) SendNotification(message string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	resp, err := c.notificationClient.Send(ctx, &pb_notification.SendRequest{Message: message})
	if err != nil {
		return "", fmt.Errorf("could not send notification: %v", err)
	}
	return resp.Result, nil
}

func (c *Client) DeleteNotification(message string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	resp, err := c.notificationClient.Delete(ctx, &pb_notification.DeleteRequest{Message: message})
	if err != nil {
		return "", fmt.Errorf("could not delete notification: %v", err)
	}
	return resp.Result, nil
}

func (c *Client) UpdateNotification(message string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	resp, err := c.notificationClient.Update(ctx, &pb_notification.UpdateRequest{Message: message})
	if err != nil {
		return "", fmt.Errorf("could not update notification: %v", err)
	}
	return resp.Result, nil
}

// Payment methods
func (c *Client) ChargePayment(message string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	resp, err := c.paymentClient.Charge(ctx, &pb_payment.ChargeRequest{Message: message})
	if err != nil {
		return "", fmt.Errorf("could not charge payment: %v", err)
	}
	return resp.Result, nil
}

func (c *Client) RefundPayment(message string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	resp, err := c.paymentClient.Refund(ctx, &pb_payment.RefundRequest{Message: message})
	if err != nil {
		return "", fmt.Errorf("could not refund payment: %v", err)
	}
	return resp.Result, nil
}

func (c *Client) GetPaymentStatus(message string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	resp, err := c.paymentClient.Status(ctx, &pb_payment.StatusRequest{Message: message})
	if err != nil {
		return "", fmt.Errorf("could not get payment status: %v", err)
	}
	return resp.Result, nil
}
