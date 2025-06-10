.PHONY: proto build run test clean

# Generate gRPC code from proto files
proto:
	protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		proto/*.proto

# Build all services
build:
	docker-compose build

# Run all services
run:
	docker-compose up -d

# Stop all services
stop:
	docker-compose down

# Run tests
test:
	go test ./...

# Clean up
clean:
	docker-compose down -v
	rm -rf proto/*.pb.go 