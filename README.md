# Plugin System with gRPC and HTTP Gateway

This project implements a configuration-driven plugin system using gRPC for service communication and HTTP Gateway for REST API access. It provides a flexible architecture for handling notifications and payments through a YAML-based plugin configuration system, all running on a single gRPC server.

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
│   └── core.go             # gRPC server and HTTP gateway setup
├── config/                 # Configuration files
│   ├── interfaces.yaml     # Interface definitions
│   └── plugins.yaml        # Plugin configurations
├── proto/                   # Protocol Buffer definitions
│   ├── notification/       # Generated notification service files
│   ├── payment/           # Generated payment service files
│   ├── google/            # Google API proto files
│   ├── notification.proto  # Notification service interface
│   └── payment.proto      # Payment service interface
├── plugins/               # Plugin implementations
│   ├── email_notification.go
│   └── stripe_payment.go
├── interfaces/           # Interface definitions and servers
│   ├── notification.go   # Notification service implementation
│   └── payment.go        # Payment service implementation
├── bin/                  # Compiled binaries
│   └── server           # Main server executable
├── main.go              # Main application entry point
├── Makefile            # Build automation
├── generate.sh         # Protocol buffer generation script
├── go.mod              # Go module definition
└── go.sum              # Go module checksums
```

## Configuration System

The project uses YAML-based configuration to define interfaces and plugins:

### Interface Configuration (`config/interfaces.yaml`)
```yaml
Interfaces:
  - name: notification
  - name: payment
```

### Plugin Configuration (`config/plugins.yaml`)
```yaml
plugins:
  - name: email
    interface: notification
    instance: NewEmailPlugin
  - name: stripe
    interface: payment
    instance: NewStripePlugin
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

4. Or run all steps at once:
   ```bash
   make all
   ```

## Available Services

All services run on a single gRPC server (Port: 50051)

### 1. Notification Service

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

### 2. Payment Service

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

The HTTP Gateway provides REST API access to all gRPC services through a single server. Example endpoints:

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
  -d '{"message": "Hello World from email notification"}'

# Delete notification
curl -X POST "http://localhost:8080/notification/delete" \
  -H "Content-Type: application/json" \
  -d '{"message": "Delete notification with ID 123"}'

# Update notification
curl -X POST "http://localhost:8080/notification/update" \
  -H "Content-Type: application/json" \
  -d '{"message": "Update notification content"}'

# Process payment charge
curl -X POST "http://localhost:8080/payment/charge" \
  -H "Content-Type: application/json" \
  -d '{"message": "Charge $100 for order #123"}'

# Process payment refund
curl -X POST "http://localhost:8080/payment/refund" \
  -H "Content-Type: application/json" \
  -d '{"message": "Refund $50 for order #123"}'

# Check payment status
curl -X POST "http://localhost:8080/payment/status" \
  -H "Content-Type: application/json" \
  -d '{"message": "Status check for transaction #456"}'
```

### gRPC Client
```go
// Example gRPC client code
conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
if err != nil {
    log.Fatal(err)
}
defer conn.Close()

// Notification service client examples
notificationClient := pb_notification.NewNotificationServiceClient(conn)

// Send notification
sendResponse, err := notificationClient.Send(context.Background(), &pb_notification.SendRequest{
    Message: "Hello World from email notification",
})

// Delete notification
deleteResponse, err := notificationClient.Delete(context.Background(), &pb_notification.DeleteRequest{
    Message: "Delete notification with ID 123",
})

// Update notification
updateResponse, err := notificationClient.Update(context.Background(), &pb_notification.UpdateRequest{
    Message: "Update notification content",
})

// Payment service client examples
paymentClient := pb_payment.NewPaymentServiceClient(conn)

// Charge payment
chargeResponse, err := paymentClient.Charge(context.Background(), &pb_payment.ChargeRequest{
    Message: "Charge $100 for order #123",
})

// Refund payment
refundResponse, err := paymentClient.Refund(context.Background(), &pb_payment.RefundRequest{
    Message: "Refund $50 for order #123",
})

// Check payment status
statusResponse, err := paymentClient.Status(context.Background(), &pb_payment.StatusRequest{
    Message: "Status check for transaction #456",
})
```

## Dependencies

- **google.golang.org/grpc** v1.70.0 - gRPC framework
- **github.com/grpc-ecosystem/grpc-gateway/v2** v2.26.3 - HTTP gateway
- **gopkg.in/yaml.v3** v3.0.1 - YAML configuration parsing
- **github.com/joho/godotenv** v1.5.1 - Environment variable loading
- **google.golang.org/protobuf** v1.36.5 - Protocol buffer support
- **google.golang.org/genproto** - Google API proto definitions

## Architecture Features

- **Configuration-Driven**: Plugins and interfaces are defined in YAML configuration files
- **Dynamic Plugin Loading**: Plugins are loaded based on configuration at runtime
- **Unified Server**: Single gRPC server handles all services with HTTP gateway
- **Interface Registry**: Dynamic interface registration system
- **Plugin Registry**: Centralized plugin management per interface type

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

4. Clean generated files:
   ```bash
   make clean
   ```

## Adding New Plugins

1. Create a new plugin implementation in the `plugins/` directory
2. Update `config/plugins.yaml` to include your new plugin
3. Ensure the interface is registered in `config/interfaces.yaml`
4. Rebuild and run the application

## Testing

The server can be tested using any gRPC client or through the HTTP Gateway endpoints. Example test requests are provided in the documentation above. The server automatically registers all configured plugins and makes them available through both gRPC and HTTP interfaces. 