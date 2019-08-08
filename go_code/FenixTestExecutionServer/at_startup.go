package FenixTestExecutionServer

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

	// Initiate internal gatewau channels
	TestExecutionGateway.InitiateGatewayChannels()

	//  Initiate the memory structure to hold all client gateway/plugin's address information
	initiateClientAddressMemoryDB()

	// Ensure that all services don't start before everything has been started
	gatewayMustStopProcessing = true

	// Initiate Database
	initiateDB(databaseFile) // If "" then Use default database file name

	// Start all Dispatch- and Transmit-Engines as a Gateway Engine and no function references, use nil
	InitiateAllTransmitAndDispatchEngines(common_code.funcTypeStruct{
		common_code.fenixOrGatewayType: common_code.GatewayEngine,
		common_code.fenixAndPluginFunctionMap: map[common_code.funcType]common_code.FuncType{
			common_code.ChannelTypeTestInstructionMessageTowardsPluginFunction:                    nil,
			common_code.ChannelTypeSupportedTestDataDomainsRequestMessageTowardsPluginFunction:    nil,
			common_code.ChannelTypeInformationMessageTowardsFenixFunction:                         nil,
			common_code.ChannelTypeTestInstructionTimeOutMessageTowardsFenixFunction:              nil,
			common_code.ChannelTypeTestExecutionLogMessageTowardsFenixFunction:                    nil,
			common_code.ChannelTypeAvailbleTestInstructionsAtPluginMessageTowardsFenixFunction:    nil,
			common_code.ChannelTypeAvailbleTestContainersAtPluginMessageTowardsFenixFunction:      nil,
			common_code.ChannelTypeTestInstructionExecutionResultMessageTowardsFenixFunction:      nil,
			common_code.ChannelTypeSupportedTestDataDomainsWithHeadersMessageTowardsFenixFunction: nil,
		},
	})

	// Listen to gRPC-calls from parent gateway/Fenix
	//startGatewayGRPCServerForMessagesTowardsFenix()

	// Listen to gRPC-calls from parent gateway/Fenix
	startGatewayGRPCServerForMessagesTowardsPlugins()

	// Listen to gRPC-calls from child gateway/plugin
	startGatewayGRPCServerForMessagesTowardsFenix()

	// Update Memory information about parent address and port with that saved in database, database overrule config-file
	updateMemoryAddressForParentAddressInfo()

	// Start all services at the same time
	gatewayMustStopProcessing = false

	// Ask clients to ReRegister them self to this gateway
	// TODO Make all Clients ReRegister them self

	// Release all saved messages to channls
	// TODO Release all massages to channels
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
