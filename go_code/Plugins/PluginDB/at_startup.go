package PluginKeyValueDBStore

// *******************************************************************
// Cleanup before closing down
//
var cleanupProcessed = false

func cleanup() {

	if cleanupProcessed == false {

		cleanupProcessed = true

		// CLose database
		closeMainDB()

		// Stop gRPC-relateds listing and close towards Plugins
		stopDbPluginGRPCServerForMessagesFromPlugins()

	}
}

// *******************************************************************
// Start all Services
//
func StartAllServices(configFileAndPath string, logfileForTest string) {
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

	// Initiate the read and write channels services
	initiateReadANdWriteServices()

	// Listen to gRPC-calls from Plugins
	startDbPluginGRPCServerForMessagesFromPlugins()

}
