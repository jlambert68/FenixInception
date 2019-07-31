package common_code

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
//gatewayClientHasBeenConnectedToParentGateway bool TODO Remove this becasue it will not be used, 190717
)

// ChannelType
const ChannelTypeInformationMessageTowardsFenix = "ChannelTypeInformationMessageTowardsFenix"
const ChannelTypeTestInstructionTimeOutMessageTowardsFenix = "ChannelTypeTestInstructionTimeOutMessageTowardsFenix"
const ChannelTypeTestExecutionLogMessageTowardsFenix = "ChannelTypeTestExecutionLogMessageTowardsFenix"
const ChannelTypeSupportedTestDataDomainsMessageTowardsFenix = "ChannelTypeSupportedTestDataDomainsMessageTowardsFenix"
const ChannelTypeAvailbleTestInstructionsAtPluginMessageTowardsFenix = "ChannelTypeAvailbleTestInstructionsAtPluginMessageTowardsFenix"
const ChannelTypeAvailbleTestContainersAtPluginMessageTowardsFenix = "ChannelTypeAvailbleTestContainersAtPluginMessageTowardsFenix"
const ChannelTypeTestInstructionExecutionResultMessageTowardsFenix = "ChannelTypeTestInstructionExecutionResultMessageTowardsFenix"
const ChannelTypeSupportedTestDataDomainsWithHeadersMessageTowardsFenix = "ChannelTypeSupportedTestDataDomainsWithHeadersMessageTowardsFenix"

// *** Internal queues used by the gateway towards Fenix ***

//  informationMessage towards Fenix
var (
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
)

var (
	registerGatewayTowardsFenixServer *grpc.Server
	gatewayTowardsFenixListener       net.Listener

	// gRPC server used to handle all traffic Towards the Plugins

)

type GRPCServerTowardsFenixStruct struct{}
