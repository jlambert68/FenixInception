package common_code

import (
	"google.golang.org/grpc"
	gRPC "jlambert/FenixInception2/go_code/TestExecutionGateway/Gateway_gRPC_api"
	"net"
)

// Data structures for clients (child-gateways and plugins) that register towards a gateway or Fenix
//TODO Remove the struct below due to it is not used, 190717

/*type gRPCClientAddressStruct struct {
	clientHasRegistered            bool
	clientIp                       string
	clientPort                     string
	clientRegistrationDateTime     string
	clientLastRegistrationDateTime string
}*/

// ChannelType
const ChannelTypeTestInstructionMessageTowardsPlugin = "ChannelTypeTestInstructionMessageTowardsPlugin"
const ChannelTypeSupportedTestDataDomainsRequestMessageTowardsPlugin = "ChannelTypeSupportedTestDataDomainsRequestMessageTowardsPlugin"

var (
	// Internal queues used by the gateway
	// TestInstruction Towards Plugin
	testInstructionMessageChannelTowardsPlugin chan *gRPC.TestInstruction_RT

	// supportedTestDataDomainsRequest Towards Plugin
	supportedTestDataDomainsRequestChannelTowardsPlugin chan *gRPC.SupportedTestDataDomainsRequest
)

var (
	registerGatewayTowardsPluginerver *grpc.Server
	gatewayTowardsPluginListener      net.Listener

	// gRPC server used to handle all traffic Towards the Plugins

)

type GRPCServerTowardsPluginStruct struct{}
