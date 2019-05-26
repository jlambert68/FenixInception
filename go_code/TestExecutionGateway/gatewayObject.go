package TestExecutionGateway

import (
	"github.com/sirupsen/logrus"
	bolt "go.etcd.io/bbolt"
)

var (
	/*
		// Connection parameters for connecting to parent Gateway/Fenix
		remoteGatewayServerConnection *grpc.ClientConn
		grpcClient                    gRPC.GatewayClient
	*/
	// Address of Parent Gateway/Fenix
	parent_address_to_dial string = ParentGatewayServer_address + ParentGatewayServer_port

	// Port where Parent Gateway/Fenix will call this gateway/client
	incomingPortForCallsFromParentGateway string
)

// The following variables is save in DB and reloaded At Startup
var (
	// Have this Gateway/Client ever been connected to parent Gateway/Fenix
	gatewayClientHasBeenConnectedToParentGateway bool
)

// Data structures for clients (child-gateways and plugins) that register towards a gateway or Fenix
type gRPCClientAddress_struct struct {
	clientHasRegistered            bool
	clientIp                       string
	clientPort                     string
	clientRegistrationDateTime     string
	clientLastRegistrationDateTime string
}

type gatewayObject_struct struct {
	// Common logger for the gateway
	logger *logrus.Logger

	// Database object used for storing any persistant data within Gateway
	db *bolt.DB

	// Database queue used for sending questions to databse
	dbMessageQueue chan dbMessage_struct
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
