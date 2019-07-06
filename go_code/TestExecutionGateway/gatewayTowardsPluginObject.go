package TestExecutionGateway

import (
	"google.golang.org/grpc"
	gRPC "jlambert/FenixInception2/go_code/TestExecutionGateway/Gateway_gRPC_api"
	"net"
)

// Data structures for clients (child-gateways and plugins) that register towards a gateway or Fenix

type gRPCClientAddressStruct struct {
	clientHasRegistered            bool
	clientIp                       string
	clientPort                     string
	clientRegistrationDateTime     string
	clientLastRegistrationDateTime string
}

// ChannelType
const channelTypeTestInstructionMessageChannel = "channelTypeTestInstructionMessageChannel"
const channelTypeSupportedTestDataDomainsRequestMessage = "channelTypeSupportedTestDataDomainsRequestMessage"

type gatewayTowardsPluginObjectStruct struct {

	// Internal queues used by the gateway
	// TestInstruction Towards Plugin
	testInstructionMessageChannel chan *gRPC.TestInstruction_RT

	// supportedTestDataDomainsRequest Towards Plugin
	supportedTestDataDomainsRequestChannel chan *gRPC.SupportedTestDataDomainsRequest
}

var (
	gatewayTowardsPluginObject        gatewayTowardsPluginObjectStruct
	registerGatewayTowardsPluginerver *grpc.Server
	gatewayTowardsPluginListener      net.Listener

	// gRPC server used to handle all traffic Towards the Plugins

)

type GRPCServerTowardsPluginStruct struct{}
