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

	// Port where Parent Gateway/Fenix will call this gateway/client
	incomingPortForCallsFromParentGateway string
)

// The following variables is saved in DB and reloaded At Startup
var (
	// Have this Gateway/Client ever been connected to parent Gateway/Fenix
	gatewayClientHasBeenConnectedToParentGateway bool
)

type gatewayTowardsFenixObjectStruct struct {

	// *** Internal queues used by the gateway ***

	//  informationMessage towards Fenix
	informationMessageChannel chan *gRPC.InformationMessage

	// testInstructionTimeOutMessage towards Fenix
	testInstructionTimeOutMessageChannel chan *gRPC.TestInstructionTimeOutMessage

	// testExecutionLogMessage towards Fenix
	testExecutionLogMessageChannel chan *gRPC.TestExecutionLogMessage

	// supportedTestDataDomainsMessage towards Fenix
	supportedTestDataDomainsMessageTowardsFenixChannel chan *gRPC.SupportedTestDataDomainsMessage

	// availbleTestInstruction<AtPluginMessage towards Fenix
	availbleTestInstructionAtPluginMessageTowardsFenixChannel chan *gRPC.AvailbleTestInstructionAtPluginMessage

	//availbleTestContainersAtPluginMessage towars Fenix
	availbleTestContainersAtPluginMessageTowardsFenixChannel chan *gRPC.AvailbleTestContainersAtPluginMessage
}

var (
	gatewayTowardsFenixObject         *gatewayTowardsFenixObjectStruct
	registerGatewayTowardsFenixServer *grpc.Server
	gatewayTowardsFenixListener       net.Listener

	// gRPC server used to handle all traffic Towards the Plugins

)

type GRPCServerTowardsFenix struct{}
