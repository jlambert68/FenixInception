# cat -e -t -v Makefile

.DEFAULT_GOAL := build

build:
	@echo "go build from proto-files..."
	protoc_3.7.0.rc2 -I Gateway_gRPC_api --go_out=plugins=grpc:Gateway_gRPC_api Gateway_gRPC_api/gateway_towards_fenix_grpc.proto
	protoc_3.7.0.rc2 -I Gateway_gRPC_api --go_out=plugins=grpc:Gateway_gRPC_api Gateway_gRPC_api/gateway_towards_plugin_grpc.proto
	protoc_3.7.0.rc2 -I Gateway_gRPC_api --go_out=plugins=grpc:Gateway_gRPC_api Gateway_gRPC_api/gateway_common_grpc.proto


check:
	cat -e -t -v Makefile