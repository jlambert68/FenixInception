package PluginKeyValueDBStore

import (
	"github.com/jlambert68/FenixInception/go_code/common_code"
	"github.com/sirupsen/logrus"
	bolt "go.etcd.io/bbolt"
)

var (
	// Common logger for the gateway
	logger *logrus.Logger

	// Database object used for storing any persistent data within Gateway
	db *bolt.DB

	// Database queue used for sending questions to database
	dbMessageQueue chan common_code.DbMessageStruct

	// Channel for informationMessage initiated in this gateway
	// removed becasue it resides in object for messages towards Fenix
	//      gatewayTowardsFenixObject.informationMessageChannelTowardsFenix chan *gRPC.InformationMessage
)

// All config parameter will be stored in the following parameter
var gatewayConfig common_code.TomlConfigStruct
