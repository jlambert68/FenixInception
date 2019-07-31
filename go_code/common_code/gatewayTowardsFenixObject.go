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

var (
	RegisterGatewayTowardsFenixServer *grpc.Server
	GatewayTowardsFenixListener       net.Listener

	// gRPC server used to handle all traffic Towards the Plugins

)

type GRPCServerTowardsFenixStruct struct{}
