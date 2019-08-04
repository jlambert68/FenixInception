package TestExecutionGateway

import (
	gRPC "github.com/jlambert68/FenixInception/go_code/common_code/Gateway_gRPC_api"
	"google.golang.org/grpc"
	"net"
)

var (
	// Connection parameters for connecting to parent Gateway/Fenix
	remoteGatewayServerConnection *grpc.ClientConn
	grpcClient                    gRPC.GatewayTowardsFenixClient

	// Port where Parent Gateway/Fenix will call this gateway/client
	//incomingPortForCallsFromParentGateway string
)

var (
	registerGatewayTowardsFenixServer *grpc.Server
	gatewayTowardsFenixListener       net.Listener

	// gRPC server used to handle all traffic Towards the Plugins

)

type gRPCServerTowardsFenixStruct struct{}
