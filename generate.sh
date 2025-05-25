#!/bin/bash

cd proto

echo "Generating code for notification service..."
protoc \
  --go_out=notification --go_opt=paths=source_relative \
  --go-grpc_out=notification --go-grpc_opt=paths=source_relative \
  --grpc-gateway_out=notification --grpc-gateway_opt=paths=source_relative \
  --openapiv2_out=notification \
  notification.proto

echo "Generating code for payment service..."
protoc \
  --go_out=payment --go_opt=paths=source_relative \
  --go-grpc_out=payment --go-grpc_opt=paths=source_relative \
  --grpc-gateway_out=payment --grpc-gateway_opt=paths=source_relative \
  --openapiv2_out=payment \
  payment.proto

echo "Code generation completed successfully."
