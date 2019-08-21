package FenixGatewayServer

import (
	"jlambert/FenixInception2/go_code/TestExecutionGateway"
)

// **********************************************************************************************************
// Save 'testInstructionExecutionResultMessage' to Main Database for Fenix Inception
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
// Save 'testInstructionExecutionResultMessage' to Main Database for Fenix Inception
//
func sendGetSupportedTestDataDomainsTowardsPlugin(testInstructionId string) (err error) {
	// Get Payload for TestInstruction to be sent to Plugin
	getSupportedTestDataDomainsMessage, err := getTestInstructionPayloadToPlugin(testInstructionId)

	if err != nil {
		// Error when loading Payload
		return err
	} else {
		// Send TestInstruction to Next Gateway
		err = dispatchEngineForTestInstructions(getSupportedTestDataDomainsMessage)

	}

	return err

}
