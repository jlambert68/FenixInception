package FenixGatewayServer

// ************************************************************************************
// Triggers functionality after the incoming message was saved in database
// In this case checking if all peers are finished and if so then trigger next peers of TestInstructions to be executed
func newIncomingTestInstructionExecutionResultMessage(peerId string) (err error) {

	var testInstructionsThatAreStillExecuting []string

	// Get all Testinstructions that are a peer to this TestInstruction, and can be run in parallell, and are still executing
	testInstructionsThatAreStillExecuting, err = listPeerTestInstructionPeersWhichIsExecuting(peerId)

	// Trigger next TestInstructions that is waiting to be executed if all current peers are finished
	if err == nil && len(testInstructionsThatAreStillExecuting) == 0 {
		testInstructionPeersThatShouldBeExecutedNext, err := listNextPeersToBeExecuted(peerId)
		if err == nil && len(testInstructionsThatAreStillExecuting) > 0 {
			err = triggerSendNextPeersForExecution(testInstructionPeersThatShouldBeExecutedNext)
		}
	}

	return err
}

// ************************************************************************************
// Triggers functionality after the incoming message was saved in database
// In this case checking if all peers are finished and if so then trigger next peers of TestInstructions to be executed
func newIncomingSupportedTestDataDomainsWithHeadersMessage() (err error) {

	err = nil
	return err
}

// ************************************************************************************
// Triggers functionality after the incoming message was saved in database
// In this case ?????????????????????????
func newIncomingTestExecutionLogMessage() (err error) {

	err = nil
	return err
}

// ************************************************************************************
// Triggers functionality after the incoming message was saved in database
// In this case ?????????????????????????
func newIncomingTestInstructionTimeOutMessage() (err error) {

	err = nil
	return err
}

// ************************************************************************************
// Triggers functionality after the incoming message was saved in database
// In this case ?????????????????????????
func newIncomingInformationMessage() (err error) {

	err = nil
	return err
}

// ************************************************************************************
// Triggers functionality after the incoming message was saved in database
// In this case ?????????????????????????
func newIncomingSupportedTestDataDomainsMessage() (err error) {

	err = nil
	return err
}

// ************************************************************************************
// Triggers functionality after the incoming message was saved in database
// In this case ?????????????????????????
func newIncomingAvailbleTestContainersAtPluginMessage() (err error) {

	err = nil
	return err
}

// ************************************************************************************
// Triggers functionality after the incoming message was saved in database
// In this case ?????????????????????????
func newIncomingAvailbleTestInstructionAtPluginMessage() (err error) {

	err = nil
	return err
}
