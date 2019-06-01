package TestExecutionGateway

import (
	"github.com/sirupsen/logrus"
	bolt "go.etcd.io/bbolt"
	gRPC "jlambert/FenixInception2/go_code/TestExecutionGateway/Gateway_gRPC_api"
)

var (
	// Common logger for the gateway
	logger *logrus.Logger

	// Database object used for storing any persistant data within Gateway
	db *bolt.DB

	// Database queue used for sending questions to databse
	dbMessageQueue chan dbMessage_struct

	// Channel for informationMessage initiated in this gateway
	localInformationMessageChannel chan *gRPC.InformationMessage
)

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

// *******************************************************************
// Information read from TOML-file
//

// All config parameter will be stored in the following parameter
var gatewayConfig tomlConfig_struct

// Main struct for toml-file
type tomlConfig_struct struct {
	gatewayIdentification gatewayIdentification_struct `toml:"gatewayIdentification"`
	systemDomain          systemDomain_struct          `toml:"systemDomain"`
	parentgRPCAddress     parentgRPCAddress_struct     `toml:"parentgRPCAddress"`
	firstClientPort       firstClientPort_struct       `toml:"firstClientPort"`
	loggingLevel          loggingLevel_struct          `toml:"loggingLevel"`
}

// local gateway information for toml-file
type gatewayIdentification_struct struct {
	callingSystemId        string
	callingSystemName      string
	callingSystemIpAddress string
}

// Ovelall domain for toml-file, Custody Cash, Compis etc.
type systemDomain_struct struct {
	callingSystemDomainId   string
	callingSystemDomainName string
}

// Parents address informaiton for toml-file
type parentgRPCAddress_struct struct {
	parentGatewayId             string
	parentGatewayServer_address string
	parentGatewayServer_port    string
}

// The first client must listen on this port, toml-file
type firstClientPort_struct struct {
	firstClientPort string
}

// The logging level that the gateway will use during runtime
type loggingLevel_struct struct {
	loggingLevel string
}

//
// *******************************************************************

// Bucket name, and some keys used in DB
const BUCKET_PARENT_ADDRESS = "Parent"
const BUCKET_KEY_PARENT_ADDRESS = "ParentId"
const BUCKET_CLIENTS = "Clients"
const BUCKET_RESEND_INFOMESSAGES_TO_FENIX = "ReSendInfoMessages"
const BUCKET_RESEND_LOG_MESSAGES_TO_FENIX = "ReSendLogMessages"
const BUCKET_RESEND_GET_TESTDATA_DOMAINS_TO_PLUGIN = "ReSendGetTestDataDomainsMessages"

// Memory Object for all clients
var clientsAddressAndPort map[string]clientsAddressAndPort_struct

// Definition for memory Object for all clients
type clientsAddressAndPort_struct struct {
	clientId      string
	clientName    string
	clientAddress string
	clientPort    string
}

// The number of messages that all internal gateway channels supports
const SUPPORTED_NUMBER_CHANNEL_MESSAGE = 100

// When gateway should start signal a INFO in log that queue in channel reached a limit
const CHANNEL_LOWER_MESSAGE_TO_BE_SIGNALED = 10

// When gateway should start signal a WARNING in log that queue in channel reached a limit
const CHANNEL_UPPER_MESSAGE_TO_BE_SIGNALED = 90
