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
	TestExecutionGateway.SetGatewayMustStopProcessingFlag = true

	// Start all Dispatch- and Transmit-Engines as a Gateway Engine inside Fenix and use function references instead of gRPC-calls
	TestExecutionGateway.InitiateAllTransmitAndDispatchEngines(common_code.FunctionsInsteadOfgRPCStruct{
		FenixOrGatewayTypeOrPlugin:                     common_code.FenixEngine,
		CallBackRegisterAvailbleTestInstructions:       CallBackRegisterAvailbleTestInstructions,
		CallBackRegistrateAailableTestContainers:       CallBackRegistrateAailableTestContainers,
		CallBackRegistrateAvailableTestDataDomains:     CallBackRegistrateAvailableTestDataDomains,
		CallBackSendTestInstructionTimeOutTowardsFenix: CallBackSendTestInstructionTimeOutTowardsFenix,
		CallBackSendTestExecutionLogTowardsFenix:       CallBackSendTestExecutionLogTowardsFenix,
		CallBackSupportedTestDataDomains:               CallBackSupportedTestDataDomains,
		CallBackSendTestInstructionResultTowardsFenix:  CallBackSendTestInstructionResultTowardsFenix,
	})

	// Start all services at the same time
	TestExecutionGateway.SetGatewayMustStopProcessingFlag = false

}

// ********************************************************************************************
// Initiate Transmit Engines for incoming messages towards Fenix and Dispatch Engines for outgoing messages towards Plugins
//

func initiateAllTransmitAndDispatchEngines() {

	// *** Towards Plugin ***
	// Start a Dispatch Engine, for 'TestInstructionMessageTowardsPlugin' as a go-routine
	go TestExecutionGateway.TransmitAndDispatchEngine(common_code.ChannelTypeTestInstructionMessageTowardsPlugin, common_code.DispatchEngineTowardsPlugin, common_code.FenixEngine, nil)

	// Start a Dispatch Engine, for 'SupportedTestDataDomainsRequestMessageTowardsPlugin,' as a go-routine
	go TestExecutionGateway.TransmitAndDispatchEngine(common_code.ChannelTypeSupportedTestDataDomainsRequestMessageTowardsPlugin, dispatchEngineTowardsPlugin, common_code.FenixEngine, nil)

	// *** Towards Fenix ***
	// Start a Transmit Engine, for 'informationMessageToBeForwarded' as a go-routine
	go TestExecutionGateway.TransmitAndDispatchEngine(common_code.ChannelTypeInformationMessageTowardsFenix, transmitEngineTowardsFenix, common_code.FenixEngine, nil)

	// Start a Transmit Engine, for 'timeOutMessageToBeForwarded' as a go-routine
	go TestExecutionGateway.TransmitAndDispatchEngine(common_code.ChannelTypeTestInstructionTimeOutMessageTowardsFenix, transmitEngineTowardsFenix, common_code.FenixEngine, nil)

	// Start a Transmit Engine, for 'spportedTestDataDomainsMessageToBeForwarded' as a go-routine
	go TestExecutionGateway.TransmitAndDispatchEngine(common_code.ChannelTypeTestExecutionLogMessageTowardsFenix, transmitEngineTowardsFenix, common_code.FenixEngine, nil)

	// Start a Transmit Engine, for 'availbleTestInstructionAtPluginMessageToBeForwarded' as a go-routine
	go TestExecutionGateway.TransmitAndDispatchEngine(common_code.ChannelTypeSupportedTestDataDomainsMessageTowardsFenix, transmitEngineTowardsFenix, common_code.FenixEngine, nil)

	// Start a Transmit Engine, for 'availbleTestContainersAtPluginMessageToBeForwarded' as a go-routine
	go TestExecutionGateway.TransmitAndDispatchEngine(common_code.ChannelTypeAvailbleTestInstructionsAtPluginMessageTowardsFenix, transmitEngineTowardsFenix, common_code.FenixEngine, nil)

	// Start a Transmit Engine, for 'availbleTestContainersAtPluginMessageToBeForwarded' as a go-routine
	go TestExecutionGateway.TransmitAndDispatchEngine(common_code.ChannelTypeAvailbleTestContainersAtPluginMessageTowardsFenix, transmitEngineTowardsFenix, common_code.FenixEngine, nil)

	// Start a Transmit Engine, for 'testInstructionExecutionResultMessageToBeForwarded' as a go-routine
	go TestExecutionGateway.TransmitAndDispatchEngine(common_code.ChannelTypeTestInstructionExecutionResultMessageTowardsFenix, transmitEngineTowardsFenix, common_code.FenixEngine, nil)

	// Start a Transmit Engine, for 'supportedTestDataDomainsWithHeadersMessageToBeForwarded' as a go-routine
	go TestExecutionGateway.TransmitAndDispatchEngine(common_code.ChannelTypeSupportedTestDataDomainsWithHeadersMessageTowardsFenix, transmitEngineTowardsFenix, common_code.FenixEngine, nil)

}
