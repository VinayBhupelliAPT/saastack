# Plugin System with gRPC

This project implements a plugin system using gRPC for communication between services.

## Prerequisites

- Go 1.16 or later
- Protocol Buffers compiler (protoc)
- Go plugins for Protocol Buffers:
  ```bash
  go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
  go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
  ```

## Building and Running

1. Generate Protocol Buffer code:
   ```bash
   make proto
   ```

2. Build the server:
   ```bash
   make build
   ```

3. Run the server:
   ```bash
   make run
   ```

## gRPC API Documentation

### Notification Service (Port: 50051)

#### Send Notification
```protobuf
rpc Send(SendRequest) returns (SendResponse)
```
- Request: `SendRequest` with message string
- Response: `SendResponse` with result string

#### Delete Notification
```protobuf
rpc Delete(DeleteRequest) returns (DeleteResponse)
```
- Request: `DeleteRequest` with message string
- Response: `DeleteResponse` with result string

#### Update Notification
```protobuf
rpc Update(UpdateRequest) returns (UpdateResponse)
```
- Request: `UpdateRequest` with message string
- Response: `UpdateResponse` with result string

### Payment Service (Port: 50052)

#### Charge Payment
```protobuf
rpc Charge(ChargeRequest) returns (ChargeResponse)
```
- Request: `ChargeRequest` with message string
- Response: `ChargeResponse` with result string

#### Refund Payment
```protobuf
rpc Refund(RefundRequest) returns (RefundResponse)
```
- Request: `RefundRequest` with message string
- Response: `RefundResponse` with result string

#### Get Payment Status
```protobuf
rpc Status(StatusRequest) returns (StatusResponse)
```
- Request: `StatusRequest` with message string
- Response: `StatusResponse` with result string

## Client Usage Example

```go
client, err := client.NewClient()
if err != nil {
    log.Fatal(err)
}

// Send notification
result, err := client.SendNotification("Hello, World!")
if err != nil {
    log.Fatal(err)
}
fmt.Println(result)

// Charge payment
result, err = client.ChargePayment("Payment for order #123")
if err != nil {
    log.Fatal(err)
}
fmt.Println(result)
```

## Project Structure

```
.
├── proto/                    # Protocol Buffer definitions
│   ├── notification.proto    # Notification service definition
│   └── payment.proto        # Payment service definition
├── server/                   # gRPC server implementation
│   └── server.go
├── client/                   # gRPC client implementation
│   └── client.go
├── plugins/                  # Plugin implementations
│   ├── email_notification.go
│   └── stripe_payment.go
├── interfaces/              # Interface definitions
│   ├── notification.go
│   └── payment.go
├── main.go                  # Main application entry point
└── Makefile                # Build automation
``` 