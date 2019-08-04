package common_code

import (
	gRPC "github.com/jlambert68/FenixInception/go_code/common_code/Gateway_gRPC_api"
	"github.com/sirupsen/logrus"
	"time"
)

/*
// ****************************************************************
// object for Integration Override of some vaiables
// Set up as default as posible
//
var GatewayInIntegrationTestMode = &gatewayIsInIntegrationTestModeStruct{
	IsInSelfIntegrationTestMode: false,
	StartWithOutAnyParent:       false,
	UsedInIntegrationTest:       false,
	ListeningOnThisPortAsParent: 0,
}
*/

//TODO Remove unused object below

/*type gatewayIsInIntegrationTestModeStruct struct {
	// Gateway is parent or child gateway as a slave for the tested gateway
	IsInSelfIntegrationTestMode bool  // Variable for handling when Gateway is in IntegrationTest-mode
	StartWithOutAnyParent       bool  // Used for Integration Tests and the gateway could be started without any parent of its own
	ListeningOnThisPortAsParent int64 // Used for Integration Tests and the gateway should used this port
	// borde alltid vara 127.0.0.1 ListeningOnThisAddressAsParent int64 // Used for Integration Tests and the gateway should used this address

	// Thia used for the tested gateway
	parentIsListeningOnThisPort    int64  // Used for Integration Tests and is used for connecting to parent
	parentIsListeningOnThisAddress string // Used for Integration Tests and is used for connecting to parent

	//
	UsedInIntegrationTest bool
}*/

// ****************************************************************

// TODO `json:"page"` fixa detta för de objekt som ska sparas i localDB

// Defines the message sent to Database Engine
type DbMessageStruct struct {
	MessageType  int                          // Will be (DbRead, DbWrite)
	Bucket       string                       // The Bucket for the message
	Key          string                       // Key to be Read or Written
	Value        []byte                       // Only used for writing messages to DB
	ResultsQueue chan<- DbResultMessageStruct // Sending function sends in which channel tp pass back the results on
}

// Used for defining Write/Read message to Database Engine
const (
	DbRead = iota
	DbWrite
	DBGetFirstObjectFromBucket
	DBDelete
)

// Message used for sending back Read-instructions from Database
type DbResultMessageStruct struct {
	Err   error  // Error message
	Key   string // Key that was Read or Written
	Value []byte // The result found in Database
}

// *******************************************************************
// Information read from TOML-file
//

// Main struct for toml-file
/*
type TomlConfigStruct struct {
	GatewayIdentification GatewayIdentificationStruct `toml:"GatewayIdentification"`
	SystemDomain          SystemDomainStruct          `toml:"SystemDomain"`
	ParentgRPCAddress     ParentgRPCAddressStruct     `toml:"ParentgRPCAddress"`
	InitialClientPort     InitialClientPortStruct     `toml:"InitialClientPort"`
	LoggingLevel          LoggingLevelStruct          `toml:"LoggingLevel"`
}
*/

type TomlConfigStruct struct {
	GatewayIdentification GatewayIdentificationStruct
	SystemDomain          SystemDomainStruct
	ParentgRPCAddress     ParentgRPCAddressStruct
	InitialClientPort     InitialClientPortStruct
	LoggingLevel          LoggingLevelStruct
	IntegrationTest       IntegrationTestStruct
}

// local gateway information for toml-file
type GatewayIdentificationStruct struct {
	GatewayId                     string
	GatewayName                   string
	GatewayIpAddress              string
	GatewaParentCallOnThisPort    int32
	GatewayChildrenCallOnThisPort int32
	CreatedDateTime               string
}

// Ovelall domain for toml-file, Custody Cash, Compis etc.
type SystemDomainStruct struct {
	GatewayDomainId   string
	GatewayDomainName string
}

// Parents address informaiton for toml-file
type ParentgRPCAddressStruct struct {
	ParentGatewayId                          string
	ParentGatewayName                        string
	ParentGatewayServerAddress               string
	ParentGatewayServerPort                  int32
	CreatedDateTime                          string
	ConnectionToParentDoneAtLeastOnce        bool
	ConnectionToParentLastConnectionDateTime string
}

// The first client, for every ip address, must listen on this port, toml-file
// Every aftercoming client , on same ip address, will add +1 to port number
type InitialClientPortStruct struct {
	InitialClientPort int32
}

// The logging level that the gateway will use during runtime
type LoggingLevelStruct struct {
	LoggingLevel logrus.Level
}

// Used for when gateway is used in Integration tests
type IntegrationTestStruct struct {
	UsedInIntegrationTest    bool
	ActAsParentGatewayInTest bool
	ActAsMainGatewayInTest   bool
	ActAsChildGatewayInTest  bool
	StartWithOutAnyParent    bool
}

//
// *******************************************************************

// Bucket name, and some keys used in DB

// Administatetive Buckets
const BucketForParentAddress = "Parent"
const BucketKeyForParentAddress = "ParentId"

//const BucketForTestInstructions = "TestInstructions" // Not used, might be used when TestInstructions should be stored for their current execution station
const BucketForGatewayIdentificationInfo = "GateWayIdentifaction"
const BucketForClients = "Clients"
const BucketKeyForGatewayIdentificationInfo = "GateWayIdentifactionId"

// Buckets used towards Fenix
const BucketForResendOfAvailableTestInstructionsTowardsFenix = "ReSendAvailableTestInstructionsTowardsFenix"
const BucketForResendOfAvailableTestContainersTowardsFenix = "ReSendAvailableTestContainersTowardsFenix"
const BucketForResendOfSupportedTestDataDomainsTowardsFenix = "ReSendSupportedTestDataDomainsTowardsFenix"
const BucketForResendOfTestInstructionExecutionResultTowardsFenix = "ReSendTestExecutionResultTowardsFenix"
const BucketForResendOfInfoMessagesTowardsFenix = "ReSendInfoMessagesTowardsFenix"
const BucketForResendOfLogMesagesTowardsFenix = "ReSendLogMessagesTowardsFenix"
const BucketForResendOTimeOutMesagesTowardsFenix = "ReSendTimeOutMessagesTowardsFenix"

// Buckets used towards Plugin
const BucketForResendTestInstructionTowardsPlugin = "ResendTestInstructionTowardsPlugin"
const BucketForResendOfGetTestdataDomainsToPlugin = "ReSendGetTestDataDomainsMessagesTowardsPlugin"

// Used for initiate TransmitDispatchEngine
const DispatchEngineTowardsPlugin = "DispatchEngineTowardsPlugin"
const TransmitEngineTowardsFenix = "TransmitEngineTowardsFenix"
const FenixEngine = "FenixEngine"
const GatewayEngine = "GatewayEngine"

// Definition for memory Object for all clients
type ClientsAddressAndPortStruct struct {
	ClientId      string
	ClientName    string
	ClientAddress string
	ClientPort    int32
}

// The number of messages that all internal gateway channels supports
const SuppertedNumberOfMessagesInChannels = 100

// When gateway should start signal a INFO in log that queue in channel reached a limit
const LowerBounderyForSignalingMessagesInChannel = 10

// When gateway should start signal a WARNING in log that queue in channel reached a limit
const UpperBounderyForSignalingMessagesInChannel = 90

// Number of seconds that Services are asleep when 'gatewayMustStopProcessing  == true'
const ServieSleepTime time.Duration = 10

// *** Internal channels used by the gateway towards Fenix ***
type GatewayChannelPackageStruct struct {

	//  informationMessage towards Fenix
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

	// Internal queues used by the gateway
	// TestInstruction Towards Plugin
	TestInstructionMessageChannelTowardsPlugin chan *gRPC.TestInstruction_RT

	// supportedTestDataDomainsRequest Towards Plugin
	SupportedTestDataDomainsRequestChannelTowardsPlugin chan *gRPC.SupportedTestDataDomainsRequest
}
