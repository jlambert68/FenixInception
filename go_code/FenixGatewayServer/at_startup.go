package FenixGatewayServer

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
	// Read 'gatewayConfig.toml' for config parameters
	processConfigFile(configFileAndPath) // Use default toml-config-file name

	// Init logger
	if logfileForTest == "" {
		//logfileForTest = generaTimeStampDateDateTime() + " - " + "localLogFile.log"
		initLogger(logfileForTest + "localLogFile.log")
	} else {
		logfileForTest = logfileForTest
		initLogger(logfileForTest)
	}

	// Initiate connection for Fenix Inceptions Main Database
	initiateMainDB()

	// Ensure that all services don't start before everything has been started
	//TODO Change this into a function-call instead
	TestExecutionGateway.SetGatewayMustStopProcessingFlag(true)

	// Start all Dispatch- and Transmit-Engines as a Gateway Engine inside Fenix and use function references instead of gRPC-calls
	TestExecutionGateway.InitiateAllTransmitAndDispatchEngines(common_code.FunctionsInsteadOfgRPCStruct{
		FenixOrGatewayTypeOrPlugin: common_code.FenixEngine,
		CallBackTowardsPlugins: common_code.CallBackTowardsPluginsType{
			CallBackSendTestInstructionTowardsPlugin:        nil,
			CallackGetSupportedTestDataDomainsTowardsPlugin: nil},
		CallBackTowardsFenix: common_code.CallBackTowardsFenixType{
			CallBackRegisterAvailbleTestInstructionsTowardsFenix:   CallBackRegisterAvailbleTestInstructions,
			CallBackRegistrateAailableTestContainersTowardsFenix:   CallBackRegistrateAailableTestContainers,
			CallBackRegistrateAvailableTestDataDomainsTowardsFenix: CallBackRegistrateAvailableTestDataDomains,
			CallBackSendMessageToFenixTowardsFenix:                 CallBackSendMessageToFenix,
			CallBackSendTestInstructionTimeOutTowardsFenix:         CallBackSendTestInstructionTimeOutTowardsFenix,
			CallBackSendTestExecutionLogTowardsFenix:               CallBackSendTestExecutionLogTowardsFenix,
			CallBackSupportedTestDataDomainsTowardsFenix:           CallBackSupportedTestDataDomains,
			CallBackSendTestInstructionResultTowardsFenixType:      CallBackSendTestInstructionResultTowardsFenix,
		},
	})

	// Start all services at the same time
	TestExecutionGateway.SetGatewayMustStopProcessingFlag(false)

}
