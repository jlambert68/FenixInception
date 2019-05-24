package TestExecutionGateway

import (
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	gRPC "jlambert/FenixInception2/go_code/TestExecutionGateway/Gateway_gRPC_api"
)

type GatewayObject_struct struct {
	logger *logrus.Logger
}

var (
	// Connection parameters for connecting to parent Gateway/Fenix
	remoteGatewayServerConnection *grpc.ClientConn
	grpcClient                    gRPC.GatewayClient

	// Address of Parent Gateway/Fenix
	parent_address_to_dial string = ParentGatewayServer_address + ParentGatewayServer_port

	// Port where Parent Gateway/Fenix will call this gateway/client
	incomingPortForCallsFromParentGateway string
)

// The following variables is save in DB and reloaded At Startup
var (
	// Have this Gateway/Client ever been connected to parent Gateway/Fenix
	gatewayClientHasBeenConnectedToParentGateway bool
)

// Data structures for clients (child-gateways and plugins) that register towards a gateway or Fenix
type gRPCClientAddress_struct struct {
	clientHasRegistered            bool
	clientIp                       string
	clientPort                     string
	clientRegistrationDateTime     string
	clientLastRegistrationDateTime string
}
