package TestExecutionGateway

import (
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	gRPC "jlambert/FenixInception2/go_code/TestExecutionGateway/Gateway_gRPC_api"
	"net"
)

type GatewayTowardsPluginObject_struct struct {

	// Common logger for the gateway
	logger *logrus.Logger

	// Internal queues used by the gateway
	// * TestInstruction Towards Plugin
	TestInstructionMessageQueue chan gRPC.TestInstruction_RT

	// gRPC server used to handle all traffic Towards the Plugins
	GRPCServer struct{}
}

/*
type clientConnectionInformation_struct struct {
	ip string
	prt string
	uuId	string
}
*/

type clientsListing_struct struct {
	clientHasRegistered            bool
	clientIp                       string
	clientPort                     string
	clientRegistrationDateTime     string
	clientLastRegistrationDateTime string
}

var (
	GatewayTowardsPluginObject        *GatewayTowardsPluginObject_struct
	RegisterGatewayTowardsPluginerver *grpc.Server
	GatewayTowardsPluginListener      net.Listener
)

/*
var (
	// Standard Worker gRPC Server
	remoteWorkerServerConnection *grpc.ClientConn
	workerClient                 worker_server_grpc_api.WorkerServerClient

	mother_address_to_dial string = common_config.MotherServer_address + common_config.MotherServer_port
)
*/

/*
// Object used in transformation from/to GUI-Object to/from DB
type GUIObjectToSave_struct struct {
	Ref         string                                    `json:"$ref"`
	Headers     []worker_server_grpc_api.HeaderStruct     `json:"headers"`
	ValueSets   []worker_server_grpc_api.ValueSetStruct   `json:"valueSets"`
	Rules       []worker_server_grpc_api.RuleStruct       `json:"rules"`
	HeaderTypes []worker_server_grpc_api.HeaderTypeStruct `json:"headerTypes"`
	RuleTypes   []worker_server_grpc_api.RuleTypeStruct   `json:"ruleTypes"`
}
*/
