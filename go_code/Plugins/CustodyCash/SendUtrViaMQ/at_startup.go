package SendUtrViaMQ

import (
	"github.com/jlambert68/FenixInception/go_code/TestExecutionGateway"
	"github.com/jlambert68/FenixInception/go_code/common_code"
)

// *******************************************************************
// Cleanup before closing down
//
var cleanupProcessed = false

func cleanup() {

	if cleanupProcessed == false {

		cleanupProcessed = true

		// CLose database
		//closeDB()

		// Stop gRPC-relateds listing and close towards Fenix
		//stopGatewayGRPCServerForMessagesTowardsFenix()

		// Stop gRPC-relateds listing and close towards Plugins
		//stopGatewayGRPCServerForMessagesTowardsPlugins()

	}
}

// *******************************************************************
// Start all Services
//
func startAllServices(configFileAndPath string, logfileForTest string, databaseFile string) {
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

	// Ensure that all services don't start before everything has been started
	//TODO Change this into a function-call instead
	//TestExecutionGateway.SetGatewayMustStopProcessingFlag(true)

	// Start all Dispatch- and Transmit-Engines as a Plugin Engine inside a Plugin and use function references instead of gRPC-calls
	gatewayOrEndpoint := common_code.FunctionsInsteadOfgRPCStruct{
		FenixOrGatewayTypeOrPlugin: common_code.PluginEngine,
		CallBackTowardsPlugins: common_code.CallBackTowardsPluginsType{
			CallBackSendTestInstructionTowardsPlugin:        CallBackSendTestInstructionTowardsPluginType,
			CallackGetSupportedTestDataDomainsTowardsPlugin: CallBackSupportedTestDataDomainsRequestTowardsPluginType},
		CallBackTowardsFenix: common_code.CallBackTowardsFenixType{
			CallBackRegisterAvailbleTestInstructionsTowardsFenix:   nil,
			CallBackRegistrateAailableTestContainersTowardsFenix:   nil,
			CallBackRegistrateAvailableTestDataDomainsTowardsFenix: nil,
			CallBackSendMessageToFenixTowardsFenix:                 nil,
			CallBackSendTestInstructionTimeOutTowardsFenix:         nil,
			CallBackSendTestExecutionLogTowardsFenix:               nil,
			CallBackSupportedTestDataDomainsTowardsFenix:           nil,
			CallBackSendTestInstructionResultTowardsFenixType:      nil,
		},
	}
	// Start TestExecutionGateway in Plugin
	TestExecutionGateway.StartAllServices(configFileAndPath, logfileForTest, databaseFile, gatewayOrEndpoint)

	// Start all services at the same time
	//TestExecutionGateway.SetGatewayMustStopProcessingFlag(false)

}
