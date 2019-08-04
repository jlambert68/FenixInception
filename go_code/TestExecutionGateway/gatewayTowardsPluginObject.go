package TestExecutionGateway

import (
	"google.golang.org/grpc"
	"net"
)

var (
	RegisterGatewayTowardsPluginerver *grpc.Server
	GatewayTowardsPluginListener      net.Listener

	// gRPC server used to handle all traffic Towards the Plugins

)
