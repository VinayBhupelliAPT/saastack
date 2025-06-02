.PHONY: all clean generate build run

all: clean generate build run

clean:
	@echo "Cleaning..."
	@rm -rf bin/
	@rm -rf interfaces/notification/proto/*.pb.go
	@rm -rf interfaces/notification/proto/*.gw.go
	@rm -rf interfaces/payment/proto/*.pb.go
	@rm -rf interfaces/payment/proto/*.gw.go

generate:
	@echo "Generating protobuf code..."
	@chmod +x ./generate.sh
	@./generate.sh

build:
	@echo "Building application..."
	@go build -o bin/server main.go

run:
	@echo "Running server..."
	@./bin/server

