package TestExecutionGateway

import (
	"github.com/sirupsen/logrus"
	bolt "go.etcd.io/bbolt"
	"google.golang.org/grpc"
	gRPC "jlambert/FenixInception2/go_code/TestExecutionGateway/Gateway_gRPC_api"
	"net"
)

type GatewayTowardsPluginObject_struct struct {

	// Common logger for the gateway
	logger *logrus.Logger

	// Database object used for storing any persistant data within Gateway
	db *bolt.DB

	// Internal queues used by the gateway
	// TestInstruction Towards Plugin
	testInstructionMessageChannel chan gRPC.TestInstruction_RT

	// supportedTestDataDomainsRequest Towards Plugin
	supportedTestDataDomainsRequestChannel chan gRPC.SupportedTestDataDomainsRequest

	//  informationMessage towards Fenix
	informationMessageChannel chan gRPC.InformationMessage

	// testInstructionTimeOutMessage towards Fenix
	testInstructionTimeOutMessageChannel chan gRPC.TestInstructionTimeOutMessage

	// testExecutionLogMessage towards Fenix
	testExecutionLogMessageTChannel chan gRPC.TestExecutionLogMessage

	// testExecutionLogMessage towards Fenix
	testExecutionLogMessage chan gRPC.TestExecutionLogMessage

	// Database queue used for sending questions to databse
	dbMessageQueue chan dbMessage_struct

	// gRPC server used to handle all traffic Towards the Plugins
	GRPCServer struct{}
}

// TODO `json:"page"` fixa detta för de objekt som ska sparas i localDB

// Defines the message sent to Database Engine
type dbMessage_struct struct {
	messageType  int                           // Will be (DB_READ, DB_WRITE)
	bucket       string                        // The Bucket for the message
	key          string                        // Key to be Read or Written
	value        []byte                        // Only used for writing messages to DB
	resultsQueue chan<- dbResultMessage_struct // Sending function sends in which channel tp pass back the results on
}

// Used for defining Write/Read message to Database Engine
const (
	DB_READ = iota
	DB_WRITE
)

// Message used for sending back Read-instructions from Database
type dbResultMessage_struct struct {
	err   error  // Error message
	key   string // Key that was Read or Written
	value []byte // The result found in Database
}

/*
type clientConnectionInformation_struct struct {
	ip string
	prt string
	uuId	string
}
*/

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
