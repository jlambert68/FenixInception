package SendUtrViaMQ

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
)

// All config parameter will be stored in the following parameter
var gatewayConfig common_code.TomlConfigStruct

// Memory Object for all clients
var clientsAddressAndPort map[string]common_code.ClientsAddressAndPortStruct

// Variable used for sync all services to be able to start and stop them at the same time
var gatewayMustStopProcessing = true

// *** Internal channels used by the gateway ***d
var gatewayChannelPackage common_code.GatewayChannelPackageStruct
