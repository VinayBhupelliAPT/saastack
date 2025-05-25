.PHONY: all clean generate build run

all: clean generate build run

clean:
	@echo "Cleaning..."
	@rm -rf bin/
	@rm -rf proto/notification/*.pb.go
	@rm -rf proto/notification/*.gw.go
	@rm -rf proto/notification/*.swagger.json
	@rm -rf proto/payment/*.pb.go
	@rm -rf proto/payment/*.gw.go
	@rm -rf proto/payment/*.swagger.json

generate:
	@echo "Generating protobuf code..."
	@chmod +x generate.sh
	@./generate.sh

build:
	@echo "Building application..."
	@go build -o bin/server main.go

run: build
	@echo "Running server..."
	@./bin/server

