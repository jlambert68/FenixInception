package PluginKeyValueDBStore

import (
	"encoding/json"
	"github.com/jlambert68/FenixInception/go_code/common_code"
	gRPC "github.com/jlambert68/FenixInception/go_code/common_code/Gateway_gRPC_api"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"strconv"
)

// *******************************************************************
// Cleanup before closing down
//
var cleanupProcessed = false

func cleanup() {

	if cleanupProcessed == false {

		cleanupProcessed = true

		// CLose database
		CloseDB()

		// Stop gRPC-relateds listing and close towards Fenix
		stopGatewayGRPCServerForMessagesTowardsFenix()

		// Stop gRPC-relateds listing and close towards Plugins
		stopGatewayGRPCServerForMessagesTowardsPlugins()

	}
}

// *******************************************************************
// Start all Services
//
func StartAllServices(configFileAndPath string, logfileForTest string, databaseFile string) {
	// Read 'pluginDbEngineConfig.toml' for config parameters
	processConfigFile(configFileAndPath) // Use default toml-config-file name

	// Init logger
	if logfileForTest == "" {
		//logfileForTest = generaTimeStampDateDateTime() + " - " + "localLogFile.log"
		initLogger(logfileForTest + "localLogFile.log")
	} else {
		logfileForTest = logfileForTest
		initLogger(logfileForTest)
	}

	// Initiate the channels used by the KeyValueStore
	initiateKeyValueStoreDbChannels()

	// Initiate the Database used by the KeyValueStore
	initiateMainDB()

	// Listen to gRPC-calls from Plugins
	startDbPluginGRPCServerForMessagesFromPlugins()

}
