package common_code

import (
	gRPC "github.com/jlambert68/FenixInception/go_code/common_code/Gateway_gRPC_api"
	"google.golang.org/grpc"
	"net"
)

var (
	// Connection parameters for connecting to parent Gateway/Fenix
	RemoteGatewayServerConnection *grpc.ClientConn
	GrpcClient                    gRPC.GatewayTowardsFenixClient

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
	InformationMessageChannelTowardsFenix chan *gRPC.InformationMessage

	// testInstructionTimeOutMessage towards Fenix
	TestInstructionTimeOutMessageChannelTowardsFenix chan *gRPC.TestInstructionTimeOutMessage

	// testExecutionLogMessage towards Fenix
	TestExecutionLogMessageChannelTowardsFenix chan *gRPC.TestExecutionLogMessage

	// supportedTestDataDomainsMessage towards Fenix
	SupportedTestDataDomainsMessageTowardsFenixChannelTowardsFenix chan *gRPC.SupportedTestDataDomainsMessage

	// availbleTestInstruction<AtPluginMessage towards Fenix
	AvailbleTestInstructionAtPluginMessageTowardsFenixChannelTowardsFenix chan *gRPC.AvailbleTestInstructionAtPluginMessage

	//availbleTestContainersAtPluginMessage towars Fenix
	AvailbleTestContainersAtPluginMessageTowardsFenixChannelTowardsFenix chan *gRPC.AvailbleTestContainersAtPluginMessage

	//availbleTestContainersAtPluginMessage towars Fenix
	TestInstructionExecutionResultMessageTowardsFenixChannelTowardsFenix chan *gRPC.TestInstructionExecutionResultMessage

	//supportedTestDataDomainsWithHeadersMessage towars Fenix
	SupportedTestDataDomainsWithHeadersMessageTowardsFenixChannelTowardsFenix chan *gRPC.SupportedTestDataDomainsWithHeadersMessage
)

var (
	RegisterGatewayTowardsFenixServer *grpc.Server
	GatewayTowardsFenixListener       net.Listener

	// gRPC server used to handle all traffic Towards the Plugins

)

type GRPCServerTowardsFenixStruct struct{}
