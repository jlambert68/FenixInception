package SendUtrViaMQ

import (
	"github.com/jlambert68/FenixInception/go_code/common_code"
	gRPC "github.com/jlambert68/FenixInception/go_code/common_code/Gateway_gRPC_api"
	"jlambert/FenixInception2/go_code/TestExecutionGateway"
)

// **********************************************************************************************************
// Send 'TestInstruction' towards Plugin
//
func sendTestInstructionTowardsPlugin(testInstructionId string) (err error) {

	err = nil

	// Get Payload for TestInstruction to be sent to Plugin
	testInstructionPayload, err := getTestInstructionPayloadToPlugin(testInstructionId)

	if err != nil {
		// Error when loading Payload
		return err
	} else {
		// Send TestInstruction to Next Gateway
		TestExecutionGateway.PutTestInstructionMessageOnChannelTowardsPlugin(testInstructionPayload)
		//err = dispatchEngineForTestInstructions(testInstructionPayload)

	}

	return err

}

// **********************************************************************************************************
// Send 'GetSuppertedTestDataDomains' towards Plugin
//
func sendGetSupportedTestDataDomainsTowardsPlugin(pluginId string, pluginName string) (err error) {

	// TODO Remove later but for now use as hardcoded when developing
	pluginId = "5df4154b-ef96-4010-a050-6cef04df55ba"
	pluginName = "Test data for Castody Cash"

	err = nil

	// Create hardcode message
	getSupportedTestDataDomainsMessage := gRPC.SupportedTestDataDomainsRequest{
		OriginalSenderId:   gatewayConfig.GatewayIdentification.GatewayId,
		OriginalSenderName: gatewayConfig.GatewayIdentification.GatewayName,
		SenderId:           gatewayConfig.GatewayIdentification.GatewayId,
		SenderName:         gatewayConfig.GatewayIdentification.GatewayName,
		PluginId:           pluginId,
		PluginName:         pluginName,
		MessageId:          common_code.GenerateUUID(logger),
		CreatedDateTime:    common_code.GeneraTimeStampUTC(),
	}
	if err != nil {
		// Error when loading Payload
		return err
	} else {
		// Send TestInstruction to Next Gateway
		TestExecutionGateway.PutGetSupportedTestDataDomainsMessageOnChannelTowardsPlugin(&getSupportedTestDataDomainsMessage)

	}

	return err

}
