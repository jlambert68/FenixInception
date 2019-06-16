# cat -e -t -v Makefile

.DEFAULT_GOAL := build

build:
	@echo "go build..."
	 go build -o gateway

run:
	./taxi_server

check:
	cat -e -t -v Makefile

proto:
	@echo "Build from proto-files..."
	protoc_3.7.0.rc2 -I Gateway_gRPC_api --go_out=plugins=grpc:Gateway_gRPC_api Gateway_gRPC_api/gateway_towards_fenix_grpc.proto
	protoc_3.7.0.rc2 -I Gateway_gRPC_api --go_out=plugins=grpc:Gateway_gRPC_api Gateway_gRPC_api/gateway_towards_plugin_grpc.proto
	protoc_3.7.0.rc2 -I Gateway_gRPC_api --go_out=plugins=grpc:Gateway_gRPC_api Gateway_gRPC_api/gateway_common_grpc.proto

test1:
	@echo "Runing IntegrationTest 1..."
	go test -tags=integration1