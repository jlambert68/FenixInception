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
	//incomingPortForCallsFromParentGateway string
)

// The following variables is saved in DB and reloaded At Startup
var (
	// Have this Gateway/Client ever been connected to parent Gateway/Fenix
	gatewayClientHasBeenConnectedToParentGateway bool
)

// ChannelType
const channelTypeInformationMessageTowardsFenix = "channelTypeinformationMessageTowardsFenix"
const channelTypeTestInstructionTimeOutMessageTowardsFenix = "channelTypeTestInstructionTimeOutMessageTowardsFenix"
const channelTypeTestExecutionLogMessageTowardsFenix = "channelTypeTestExecutionLogMessageTowardsFenix"
const channelTypeSupportedTestDataDomainsMessageTowardsFenix = "channelTypeSupportedTestDataDomainsMessageTowardsFenix"
const channelTypeAvailbleTestInstructionsAtPluginMessageTowardsFenix = "channelTypeAvailbleTestInstructionsAtPluginMessageTowardsFenix"
const channelTypeAvailbleTestContainersAtPluginMessageTowardsFenix = "channelTypeAvailbleTestContainersAtPluginMessageTowardsFenix"
const channelTypeTestInstructionExecutionResultMessageTowardsFenix = "channelTypeTestInstructionExecutionResultMessageTowardsFenix"
const channelTypeSupportedTestDataDomainsWithHeadersMessageTowardsFenix = "channelTypeSupportedTestDataDomainsWithHeadersMessageTowardsFenix"

type gatewayTowardsFenixObjectStruct struct {
	// *** Internal queues used by the gateway ***

	//  informationMessage towards Fenix

	informationMessageChannelTowardsFenix chan *gRPC.InformationMessage

	// testInstructionTimeOutMessage towards Fenix
	testInstructionTimeOutMessageChannelTowardsFenix chan *gRPC.TestInstructionTimeOutMessage

	// testExecutionLogMessage towards Fenix
	testExecutionLogMessageChannelTowardsFenix chan *gRPC.TestExecutionLogMessage

	// supportedTestDataDomainsMessage towards Fenix
	supportedTestDataDomainsMessageTowardsFenixChannelTowardsFenix chan *gRPC.SupportedTestDataDomainsMessage

	// availbleTestInstruction<AtPluginMessage towards Fenix
	availbleTestInstructionAtPluginMessageTowardsFenixChannelTowardsFenix chan *gRPC.AvailbleTestInstructionAtPluginMessage

	//availbleTestContainersAtPluginMessage towars Fenix
	availbleTestContainersAtPluginMessageTowardsFenixChannelTowardsFenix chan *gRPC.AvailbleTestContainersAtPluginMessage

	//availbleTestContainersAtPluginMessage towars Fenix
	testInstructionExecutionResultMessageTowardsFenixChannelTowardsFenix chan *gRPC.TestInstructionExecutionResultMessage

	//supportedTestDataDomainsWithHeadersMessage towars Fenix
	supportedTestDataDomainsWithHeadersMessageTowardsFenixChannelTowardsFenix chan *gRPC.SupportedTestDataDomainsWithHeadersMessage
}

var (
	gatewayTowardsFenixObject         gatewayTowardsFenixObjectStruct
	registerGatewayTowardsFenixServer *grpc.Server
	gatewayTowardsFenixListener       net.Listener

	// gRPC server used to handle all traffic Towards the Plugins

)

type GRPCServerTowardsFenixStruct struct{}
