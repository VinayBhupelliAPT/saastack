#!/bin/bash

echo "Generating code for notification service..."
protoc \
  -Iinterfaces \
  --go_out=interfaces/notification/proto --go_opt=paths=import \
  --go-grpc_out=interfaces/notification/proto --go-grpc_opt=paths=import \
  --grpc-gateway_out=interfaces/notification/proto --grpc-gateway_opt=paths=import \
  interfaces/notification/notification.proto

echo "Generating code for payment service..."
protoc \
  -Iinterfaces \
  --go_out=interfaces/payment/proto --go_opt=paths=import \
  --go-grpc_out=interfaces/payment/proto --go-grpc_opt=paths=import \
  --grpc-gateway_out=interfaces/payment/proto --grpc-gateway_opt=paths=import \
  interfaces/payment/payment.proto


echo "Generating code for bookstore service..."
protoc \
  -Iinterfaces \
  --go_out=interfaces/bookstore/proto --go_opt=paths=import \
  --go-grpc_out=interfaces/bookstore/proto --go-grpc_opt=paths=import \
  --grpc-gateway_out=interfaces/bookstore/proto --grpc-gateway_opt=paths=import \
  interfaces/bookstore/bookstore.proto

echo "Code generation completed successfully."
