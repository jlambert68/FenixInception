package TestExecutionGateway

import (
	"google.golang.org/grpc"
	"net"
)

var (
	registerGatewayTowardsPluginerver *grpc.Server
	gatewayTowardsPluginListener      net.Listener

	// gRPC server used to handle all traffic Towards the Plugins

)

type gRPCServerTowardsPluginStruct struct{}
