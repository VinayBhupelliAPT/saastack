# Plugin System with gRPC and HTTP Gateway

This project implements a plugin system using gRPC for service communication and HTTP Gateway for REST API access. It provides a flexible architecture for handling notifications and payments through a plugin-based system.

## Prerequisites

- Go 1.23 or later
- Protocol Buffers compiler (protoc)
- Go plugins for Protocol Buffers:
  ```bash
  go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
  go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
  go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest
  ```

## Project Structure

```
.
├── core/                    # Core plugin system implementation
│   └── core.go             # Plugin registry and routing logic
├── proto/                   # Protocol Buffer definitions
│   ├── notification/       # Notification service definitions
│   ├── payment/           # Payment service definitions
│   ├── notification.proto  # Notification service interface
│   └── payment.proto      # Payment service interface
├── server/                 # gRPC and HTTP server implementation
│   └── server.go
├── plugins/               # Plugin implementations
│   ├── email_notification.go
│   └── stripe_payment.go
├── interfaces/           # Interface definitions
│   ├── notification.go
│   └── payment.go
├── main.go              # Main application entry point
├── Makefile            # Build automation
└── generate.sh         # Protocol buffer generation script
```

## Building and Running

1. Generate Protocol Buffer code:
   ```bash
   make generate
   ```

2. Build the server:
   ```bash
   make build
   ```

3. Run the server:
   ```bash
   make run
   ```

## Available Services

### 1. Notification Service (Port: 50051)

#### Methods:
- **Send**: Send a notification message
  ```protobuf
  rpc Send(SendRequest) returns (SendResponse)
  ```
- **Delete**: Delete a notification
  ```protobuf
  rpc Delete(DeleteRequest) returns (DeleteResponse)
  ```
- **Update**: Update a notification
  ```protobuf
  rpc Update(UpdateRequest) returns (UpdateResponse)
  ```

### 2. Payment Service (Port: 50052)

#### Methods:
- **Charge**: Process a payment
  ```protobuf
  rpc Charge(ChargeRequest) returns (ChargeResponse)
  ```
- **Refund**: Process a refund
  ```protobuf
  rpc Refund(RefundRequest) returns (RefundResponse)
  ```
- **Status**: Check payment status
  ```protobuf
  rpc Status(StatusRequest) returns (StatusResponse)
  ```

## HTTP Gateway (Port: 8080)

The HTTP Gateway provides REST API access to all gRPC services. Example endpoints:

### Notification Endpoints
```
POST /notification/send
POST /notification/delete
POST /notification/update
```

### Payment Endpoints
```
POST /payment/charge
POST /payment/refund
POST /payment/status
```

## Example Usage

### HTTP API
```bash
# Send notification
curl -X POST "http://localhost:8080/notification/send" \
  -H "Content-Type: application/json" \
  -d '{"message": "Hello World", "plugin": "email"}'

# Process payment
curl -X POST "http://localhost:8080/payment/charge" \
  -H "Content-Type: application/json" \
  -d '{"message": "Payment for order #123", "plugin": "stripe"}'
```

### gRPC Client
```go
// Example gRPC client code
conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
if err != nil {
    log.Fatal(err)
}
defer conn.Close()

client := pb_notification.NewNotificationServiceClient(conn)
response, err := client.Send(context.Background(), &pb_notification.SendRequest{
    Message: "Hello World",
    Plugin:  "email",
})
```

## Dependencies

- google.golang.org/grpc v1.70.0
- github.com/grpc-ecosystem/grpc-gateway/v2 v2.26.3
- gopkg.in/yaml.v3 v3.0.1
- google.golang.org/protobuf v1.36.5

## Development

1. Install dependencies:
   ```bash
   go mod download
   ```

2. Generate protocol buffer code:
   ```bash
   make generate
   ```

3. Build and run:
   ```bash
   make run
   ```

## Testing

The server can be tested using any gRPC client or through the HTTP Gateway endpoints. Example test requests are provided in the documentation above. 