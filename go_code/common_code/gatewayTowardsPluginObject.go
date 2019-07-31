package common_code

import (
	gRPC "github.com/jlambert68/FenixInception/go_code/common_code/Gateway_gRPC_api"
	"google.golang.org/grpc"
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
	TestInstructionMessageChannelTowardsPlugin chan *gRPC.TestInstruction_RT

	// supportedTestDataDomainsRequest Towards Plugin
	SupportedTestDataDomainsRequestChannelTowardsPlugin chan *gRPC.SupportedTestDataDomainsRequest
)

var (
	RegisterGatewayTowardsPluginerver *grpc.Server
	GatewayTowardsPluginListener      net.Listener

	// gRPC server used to handle all traffic Towards the Plugins

)

type GRPCServerTowardsPluginStruct struct{}
