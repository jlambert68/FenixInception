package FenixGatewayServer

import (
	"google.golang.org/grpc"
	"net"
)

var (
	registerFenixGatewayTowardsPluginerver *grpc.Server
	fenixGatewayTowardsPluginListener      net.Listener

	// gRPC server used to handle all traffic Towards the Plugins

)

type gRPCServerTowardsPluginStruct struct{}
