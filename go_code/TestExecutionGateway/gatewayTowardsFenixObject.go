package TestExecutionGateway

import (
	"google.golang.org/grpc"
	gRPC "jlambert/FenixInception2/go_code/TestExecutionGateway/Gateway_gRPC_api"
	"net"
)

var (

	// Connection parameters for connecting to parent Gateway/Fenix
	remoteGatewayServerConnection *grpc.ClientConn
	grpcClient                    gRPC.GatewayTowardsFenixClient

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

type gatewayTowardsFenixObject_struct struct {

	// Inherit common objects
	gatewayCommonObjects *gatewayObject_struct

	// *** Internal queues used by the gateway ***

	//  informationMessage towards Fenix
	informationMessageChannel chan *gRPC.InformationMessage

	// testInstructionTimeOutMessage towards Fenix
	testInstructionTimeOutMessageChannel chan *gRPC.TestInstructionTimeOutMessage

	// testExecutionLogMessage towards Fenix
	testExecutionLogMessageChannel chan *gRPC.TestExecutionLogMessage

	// supportedTestDataDomainsMessage towards Fenix
	supportedTestDataDomainsMessageTowardsFenixChannel chan *gRPC.SupportedTestDataDomainsMessage
}

var (
	gatewayTowardsFenixObject         *gatewayTowardsPluginObject_struct
	registerGatewayTowardsFenixServer *grpc.Server
	gatewayTowardsFenixListener       net.Listener

	// gRPC server used to handle all traffic Towards the Plugins

)

type GRPCServerTowardsFenix struct{}
