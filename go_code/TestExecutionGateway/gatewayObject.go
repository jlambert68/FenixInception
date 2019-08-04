package TestExecutionGateway

import (
	"github.com/jlambert68/FenixInception/go_code/common_code"
	"github.com/sirupsen/logrus"
	bolt "go.etcd.io/bbolt"
)

var (
	// Common logger for the gateway
	Logger *logrus.Logger

	// Database object used for storing any persistent data within Gateway
	db *bolt.DB

	// Database queue used for sending questions to database
	DbMessageQueue chan common_code.DbMessageStruct

	// Channel for informationMessage initiated in this gateway
	// removed becasue it resides in object for messages towards Fenix
	//      gatewayTowardsFenixObject.informationMessageChannelTowardsFenix chan *gRPC.InformationMessage
)

// All config parameter will be stored in the following parameter
var GatewayConfig common_code.TomlConfigStruct

// Memory Object for all clients
var ClientsAddressAndPort map[string]common_code.ClientsAddressAndPortStruct

// Variable used for sync all services to be able to start and stop them at the same time
var GatewayMustStopProcessing = true
