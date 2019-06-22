package TestExecutionGateway

import (
	"github.com/sirupsen/logrus"
	bolt "go.etcd.io/bbolt"
	gRPC "jlambert/FenixInception2/go_code/TestExecutionGateway/Gateway_gRPC_api"
	"time"
)

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

type gatewayIsInIntegrationTestModeStruct struct {
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
}

// ****************************************************************

// TODO fortsätt här

var (
	// Common logger for the gateway
	logger *logrus.Logger

	// Database object used for storing any persistant data within Gateway
	db *bolt.DB

	// Database queue used for sending questions to databse
	dbMessageQueue chan dbMessageStruct

	// Channel for informationMessage initiated in this gateway
	localInformationMessageChannel chan *gRPC.InformationMessage
)

// TODO `json:"page"` fixa detta för de objekt som ska sparas i localDB

// Defines the message sent to Database Engine
type dbMessageStruct struct {
	messageType  int                          // Will be (DbRead, DbWrite)
	bucket       string                       // The Bucket for the message
	key          string                       // Key to be Read or Written
	value        []byte                       // Only used for writing messages to DB
	resultsQueue chan<- dbResultMessageStruct // Sending function sends in which channel tp pass back the results on
}

// Used for defining Write/Read message to Database Engine
const (
	DbRead = iota
	DbWrite
)

// Message used for sending back Read-instructions from Database
type dbResultMessageStruct struct {
	err   error  // Error message
	key   string // Key that was Read or Written
	value []byte // The result found in Database
}

// *******************************************************************
// Information read from TOML-file
//

// All config parameter will be stored in the following parameter
var gatewayConfig TomlConfigStruct

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
const BucketForParentAddress = "Parent"
const BucketKeyForParentAddress = "ParentId"
const BucketForClients = "Clients"
const BucketForResendOfInfoMessagesToFenix = "ReSendInfoMessages"
const BucketForResendOfLogMesagesToFenix = "ReSendLogMessages"
const BucketForResendOfGetTestdataDomainsToPlugin = "ReSendGetTestDataDomainsMessages"
const BucketForTestInstructions = "TestInstructions"
const BucketForGatewayIdentificationInfo = "GateWayIdentifaction"
const BucketKeyForGatewayIdentificationInfo = "GateWayIdentifactionId"
const BucketForResendOfAvailableTestInstructionsToFenix = "ReSendAvailableTestInstructions"
const BucketForResendOfAvailableTestContainers = "ReSendAvailableTestContainers"
const BucketForResendOfSupportedTestDataDomains = "ReSendSupportedTestDataDomains"
const BucketForResendOfTestInstructionExecutionResult = "ReSendTestExecutionResult"

// Memory Object for all clients
var clientsAddressAndPort map[string]clientsAddressAndPortStruct

// Definition for memory Object for all clients
type clientsAddressAndPortStruct struct {
	clientId      string
	clientName    string
	clientAddress string
	clientPort    int32
}

// The number of messages that all internal gateway channels supports
const SuppertedNumberOfMessagesInChannels = 100

// When gateway should start signal a INFO in log that queue in channel reached a limit
const LowerBounderyForSignalingMessagesInChannel = 10

// When gateway should start signal a WARNING in log that queue in channel reached a limit
const UpperBounderyForSignalingMessagesInChannel = 90

// Variable used for sync all services to be able to start and stop them at the same time
var gatewayMustStopProcessing = true

// Number of seconds that Services are asleep when 'gatewayMustStopProcessing  == true'
const ServieSleepTime time.Duration = 10
